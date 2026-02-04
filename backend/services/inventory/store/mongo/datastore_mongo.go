// Copyright 2023 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package mongo

import (
	"context"
	"crypto/tls"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	mopts "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-server/pkg/identity"
	"github.com/mendersoftware/mender-server/pkg/log"
	"github.com/mendersoftware/mender-server/pkg/tiers"

	"github.com/mendersoftware/mender-server/services/inventory/model"
	"github.com/mendersoftware/mender-server/services/inventory/store"
	"github.com/mendersoftware/mender-server/services/inventory/utils"
)

const (
	DbVersion = "1.1.0"

	DbName        = "inventory"
	DbDevicesColl = "devices"

	DBKeyTenantID = "tenant_id"
	DBKeyID       = "_id"

	DbDevId              = "_id"
	DbDevAttributes      = "attributes"
	DbDevGroup           = "group"
	DbDevRevision        = "revision"
	DbDevUpdatedTs       = "updated_ts"
	DbDevAttributesText  = "text"
	DbDevAttributesTs    = "timestamp"
	DbDevAttributesDesc  = "description"
	DbDevAttributesValue = "value"
	DbDevAttributesScope = "scope"
	DbDevAttributesName  = "name"

	DbScopeInventory = "inventory"

	FiltersAttributesMaxDevices = 5000
	FiltersAttributesLimit      = 500
)

var (
	//with offcial mongodb supported driver we keep client
	clientGlobal *mongo.Client

	// once ensures client is created only once
	once sync.Once

	ErrNotFound = errors.New("mongo: no documents in result")
)

type DataStoreMongoConfig struct {
	// connection string
	ConnectionString string

	// SSL support
	SSL           bool
	SSLSkipVerify bool

	// Overwrites credentials provided in connection string if provided
	Username string
	Password string
}

type DataStoreMongo struct {
	client      *mongo.Client
	automigrate bool
}

func NewDataStoreMongoWithSession(client *mongo.Client) store.DataStore {
	return &DataStoreMongo{client: client}
}

// config.ConnectionString must contain a valid
func NewDataStoreMongo(config DataStoreMongoConfig) (store.DataStore, error) {
	//init master session
	var err error
	once.Do(func() {
		if !strings.Contains(config.ConnectionString, "://") {
			config.ConnectionString = "mongodb://" + config.ConnectionString
		}
		clientOptions := mopts.Client().ApplyURI(config.ConnectionString)

		if config.Username != "" {
			clientOptions.SetAuth(mopts.Credential{
				Username: config.Username,
				Password: config.Password,
			})
		}

		if config.SSL {
			tlsConfig := &tls.Config{}
			tlsConfig.InsecureSkipVerify = config.SSLSkipVerify
			clientOptions.SetTLSConfig(tlsConfig)
		}
		reg := bson.NewRegistry()
		reg.RegisterTypeDecoder(tDeviceP, bsoncodec.ValueDecoderFunc(deviceDecodeValue))
		reg.RegisterTypeDecoder(tDevice, bsoncodec.ValueDecoderFunc(deviceDecodeValue))
		clientOptions.SetRegistry(reg)
		log.NewEmpty().Warn("///////", reg)

		ctx := context.Background()
		l := log.FromContext(ctx)
		clientGlobal, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			l.Errorf("mongo: error connecting to mongo '%s'", err.Error())
			return
		}
		if clientGlobal == nil {
			l.Errorf("mongo: client is nil. wow.")
			return
		}
		// from: https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
		/*
			It is best practice to keep a client that is connected to MongoDB around so that the
			application can make use of connection pooling - you don't want to open and close a
			connection for each query. However, if your application no longer requires a connection,
			the connection can be closed with client.Disconnect() like so:
		*/
		err = clientGlobal.Ping(ctx, nil)
		if err != nil {
			clientGlobal = nil
			l.Errorf("mongo: error pinging mongo '%s'", err.Error())
			return
		}
		if clientGlobal == nil {
			l.Errorf("mongo: global instance of client is nil.")
			return
		}
	})

	if clientGlobal == nil {
		return nil, errors.New("failed to open mongo-driver session")
	}
	db := &DataStoreMongo{client: clientGlobal}

	return db, nil
}

func (db *DataStoreMongo) Ping(ctx context.Context) error {
	res := db.client.Database(DbName).RunCommand(ctx, bson.M{"ping": 1})
	return res.Err()
}

func tenantIDFromContext(ctx context.Context) (tenantID string) {
	if id := identity.FromContext(ctx); id != nil {
		tenantID = id.Tenant
	}
	return tenantID
}

func (db *DataStoreMongo) GetDevices(
	ctx context.Context,
	q store.ListQuery,
) ([]model.Device, int, error) {
	c := db.client.Database(DbName).Collection(DbDevicesColl)

	findQuery := q.Filters.ToMongoFilter()
	findQuery[DBKeyTenantID] = tenantIDFromContext(ctx)
	if q.GroupName != "" {
		findQuery["group"] = q.GroupName
	} else if q.HasGroup != nil {
		findQuery["group"] = bson.M{
			"$exists": *q.HasGroup,
		}
	}

	findOptions := mopts.Find()
	if q.Skip > 0 {
		findOptions.SetSkip(int64(q.Skip))
	}
	if q.Limit > 0 {
		findOptions.SetLimit(int64(q.Limit))
	}
	// FIXME:
	// if q.Sort != nil {
	// 	name := fmt.Sprintf(
	// 		"%s-%s",
	// 		q.Sort.AttrScope,
	// 		model.GetDeviceAttributeNameReplacer().Replace(q.Sort.AttrName),
	// 	)
	// 	sortField := fmt.Sprintf("%s.%s.%s", DbDevAttributes, name, DbDevAttributesValue)
	// 	sortFieldQuery := bson.D{{Key: sortField, Value: 1}}
	// 	if !q.Sort.Ascending {
	// 		sortFieldQuery[0].Value = -1
	// 	}
	// 	findOptions.SetSort(sortFieldQuery)
	// }

	cursor, err := c.Find(ctx, findQuery, findOptions)
	if err != nil {
		return nil, -1, errors.Wrap(err, "failed to search devices")
	}
	defer cursor.Close(ctx)

	devices := []model.Device{}
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, -1, errors.Wrap(err, "failed to search devices")
	}

	count, err := c.CountDocuments(ctx, findQuery)
	if err != nil {
		return nil, -1, errors.Wrap(err, "failed to count devices")
	}

	return devices, int(count), nil
}

func (db *DataStoreMongo) GetDevice(
	ctx context.Context,
	id model.DeviceID,
) (*model.Device, error) {
	var res model.Device
	c := db.client.
		Database(DbName).
		Collection(DbDevicesColl)
	l := log.FromContext(ctx)

	if id == model.NilDeviceID {
		return nil, nil
	}
	if err := c.FindOne(ctx, bson.M{
		DbDevId:       id,
		DBKeyTenantID: tenantIDFromContext(ctx)},
	).
		Decode(&res); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, nil
		default:
			l.Errorf("GetDevice: %v", err)
			return nil, errors.Wrap(err, "failed to fetch device")
		}
	}
	return &res, nil
}

// AddDevice inserts a new device, initializing the inventory data.
func (db *DataStoreMongo) AddDevice(ctx context.Context, dev *model.Device) error {
	_, err := db.UpsertDevicesAttributesWithUpdated(
		ctx, []model.DeviceID{dev.ID}, dev.Attributes, "", "",
	)
	if err != nil {
		return errors.Wrap(err, "failed to store device")
	}
	return nil
}

func (db *DataStoreMongo) UpsertDevicesAttributesWithRevision(
	ctx context.Context,
	devices []model.DeviceUpdate,
	attrs model.DeviceAttributes,
) (*model.UpdateResult, error) {
	return db.upsertAttributes(ctx, devices, attrs, false, true, "", "", nil)
}

func (db *DataStoreMongo) UpsertDevicesAttributesWithUpdated(
	ctx context.Context,
	ids []model.DeviceID,
	attrs model.DeviceAttributes,
	scope string,
	etag string,
) (*model.UpdateResult, error) {
	withUpdated := scope == model.AttrScopeInventory
	return db.upsertAttributes(
		ctx, makeDevsWithIds(ids),
		attrs, withUpdated,
		false, scope,
		etag, nil,
	)
}

func (db *DataStoreMongo) UpsertDevicesAttributes(
	ctx context.Context,
	ids []model.DeviceID,
	attrs model.DeviceAttributes,
	notModifiedAfter *time.Time,
) (*model.UpdateResult, error) {
	return db.upsertAttributes(
		ctx, makeDevsWithIds(ids), attrs,
		false, false, "",
		"", notModifiedAfter,
	)
}

func makeDevsWithIds(ids []model.DeviceID) []model.DeviceUpdate {
	devices := make([]model.DeviceUpdate, len(ids))
	for i, id := range ids {
		devices[i].Id = id
	}
	return devices
}

type DeviceAttributes model.DeviceAttributes

func (attrs DeviceAttributes) MarshalBSON() ([]byte, error) {
	ret := bson.M{}
	for _, attr := range attrs {
		scope := attr.Scope
		switch scope {
		case model.AttrScopeIdentity:
			if attr.Name == "status" {
				ret[attr.Name] = attr.Value
			} else {
				if a, ok := ret[scope].([]model.DeviceAttribute); ok {
					ret[scope] = append(a, attr)
				} else {
					ret[scope] = []model.DeviceAttribute{attr}
				}
			}

		case model.AttrScopeInventory, "":
			scope = model.AttrScopeInventory
			fallthrough
		case model.AttrScopeTags,
			model.AttrScopeMonitor:
			if a, ok := ret[scope].([]model.DeviceAttribute); ok {
				ret[scope] = append(a, attr)
			} else {
				ret[scope] = []model.DeviceAttribute{attr}
			}

		case model.AttrScopeSystem:
			ret[attr.Name] = attr.Value

		default:
			return nil, fmt.Errorf("unknown scope %q", attr.Scope)
		}
	}
	return bson.Marshal(ret)
}

func (db *DataStoreMongo) upsertAttributes(
	ctx context.Context,
	devices []model.DeviceUpdate,
	attrs model.DeviceAttributes,
	withUpdated bool,
	withRevision bool,
	scope string,
	etag string,
	notModifiedAfter *time.Time,
) (*model.UpdateResult, error) {
	const createdField = model.AttrNameCreated
	const etagField = model.AttrNameTagsEtag
	const updatedTS = DbDevUpdatedTs
	var (
		result *model.UpdateResult
		err    error
	)

	c := db.client.
		Database(DbName).
		Collection(DbDevicesColl)

	now := time.Now()
	tenantID := tenantIDFromContext(ctx)
	oninsert := bson.M{
		createdField:  now,
		DBKeyTenantID: tenantID,
	}
	updateAttrs := DeviceAttributes(attrs)
	if !withRevision {
		updateAttrs = append(updateAttrs, model.DeviceAttribute{
			Name:  DbDevRevision,
			Value: 0,
			Scope: model.AttrScopeSystem,
		})
	}

	if withUpdated {
		const updatedField = model.AttrNameUpdated
		updateAttrs = append(updateAttrs, model.DeviceAttribute{
			Name:  updatedField,
			Value: now,
			Scope: model.AttrScopeSystem,
		})
	}

	switch len(devices) {
	case 0:
		return &model.UpdateResult{}, nil
	case 1:
		filter := bson.M{
			"_id":         devices[0].Id,
			DBKeyTenantID: tenantID,
		}
		updateOpts := mopts.FindOneAndUpdate().
			SetUpsert(true).
			SetReturnDocument(mopts.After)

		if withRevision {
			filter[DbDevRevision] = bson.M{"$lt": devices[0].Revision}
			updateAttrs = append(updateAttrs, model.DeviceAttribute{
				Name:  DbDevRevision,
				Value: devices[0].Revision,
				Scope: model.AttrScopeSystem,
			})
		}
		if scope == model.AttrScopeTags {
			updateAttrs = append(updateAttrs, model.DeviceAttribute{
				Name:  DbDevRevision,
				Value: uuid.New().String(),
				Scope: model.AttrScopeSystem,
			})
			updateOpts = mopts.FindOneAndUpdate().
				SetUpsert(false).
				SetReturnDocument(mopts.After)
		}
		if etag != "" {
			filter[etagField] = bson.M{"$eq": etag}
		}
		if notModifiedAfter != nil {
			filter["$or"] = []bson.M{
				{updatedTS: bson.M{"$lte": notModifiedAfter}},
				{updatedTS: bson.M{"$exists": false}},
			}
		}

		update := bson.M{
			"$set":         updateAttrs,
			"$setOnInsert": oninsert,
		}

		device := &model.Device{}
		res := c.FindOneAndUpdate(ctx, filter, update, updateOpts)
		err = res.Decode(device)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, store.ErrWriteConflict
			} else if err == mongo.ErrNoDocuments {
				return &model.UpdateResult{}, nil
			} else {
				return nil, err
			}
		}
		result = &model.UpdateResult{
			MatchedCount: 1,
			CreatedCount: 0,
			Devices:      []*model.Device{device},
		}
	default:
		var bres *mongo.BulkWriteResult
		// Perform single bulk-write operation
		// NOTE: Can't use UpdateMany as $in query operator does not
		//       upsert missing devices.

		models := make([]mongo.WriteModel, len(devices))
		for i, dev := range devices {
			umod := mongo.NewUpdateOneModel()
			filter := bson.M{"_id": dev.Id, DBKeyTenantID: tenantID}
			update := slices.Clone(updateAttrs)
			if withRevision {
				filter[DbDevRevision] = bson.M{"$lt": dev.Revision}
				update = append(update, model.DeviceAttribute{
					Name:  DbDevRevision,
					Value: dev.Revision,
					Scope: model.AttrScopeSystem,
				})
			}
			if notModifiedAfter != nil {
				filter["$or"] = []bson.M{
					{updatedTS: bson.M{"$lte": notModifiedAfter}},
					{updatedTS: bson.M{"$exists": false}},
				}
			}
			umod.Update = bson.M{
				"$set":         update,
				"$setOnInsert": oninsert,
			}
			umod.Filter = filter
			umod.SetUpsert(true)
			models[i] = umod
		}
		bres, err = c.BulkWrite(
			ctx, models, mopts.BulkWrite().SetOrdered(false),
		)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				// bulk mode, swallow the error as we already updated the other devices
				// and the Matchedcount and CreatedCount values will tell the caller if
				// all the operations succeeded or not
				err = nil
			} else {
				return nil, err
			}
		}
		result = &model.UpdateResult{
			MatchedCount: bres.MatchedCount,
			CreatedCount: bres.UpsertedCount,
		}
	}
	return result, err
}

func (db *DataStoreMongo) UpsertRemoveDeviceAttributes(
	ctx context.Context,
	id model.DeviceID,
	updateAttrs model.DeviceAttributes,
	removeAttrs model.DeviceAttributes,
	scope string,
	etag string,
) (*model.UpdateResult, error) {
	const updatedField = model.AttrNameUpdated
	const createdField = model.AttrNameCreated
	const etagField = model.AttrNameTagsEtag
	var (
		err error
	)
	now := time.Now()

	c := db.client.
		Database(DbName).
		Collection(DbDevicesColl)

	filter := bson.M{DBKeyID: id, DBKeyTenantID: tenantIDFromContext(ctx)}
	if etag != "" {
		filter[etagField] = bson.M{"$eq": etag}
	}

	updateOpts := mopts.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(mopts.After)
	if scope == model.AttrScopeTags {
		updateAttrs = append(updateAttrs, model.DeviceAttribute{
			Name:  etagField,
			Value: uuid.New().String(),
			Scope: model.AttrScopeSystem,
		})
		updateOpts = updateOpts.SetUpsert(false)
	}
	if scope != model.AttrScopeTags {
		updateAttrs = append(updateAttrs, model.DeviceAttribute{
			Name:  model.AttrNameUpdated,
			Value: now,
			Scope: model.AttrScopeSystem,
		})
	}
	update := bson.M{
		"$set": DeviceAttributes(updateAttrs),
		"$setOnInsert": bson.M{
			createdField:  now,
			DBKeyTenantID: tenantIDFromContext(ctx),
		},
	}
	device := &model.Device{}
	res := c.FindOneAndUpdate(ctx, filter, update, updateOpts)
	err = res.Decode(device)
	if err == mongo.ErrNoDocuments {
		return &model.UpdateResult{
			MatchedCount: 0,
			CreatedCount: 0,
			Devices:      []*model.Device{},
		}, nil
	} else if err == nil {
		return &model.UpdateResult{
			MatchedCount: 1,
			CreatedCount: 0,
			Devices:      []*model.Device{device},
		}, nil
	}
	return nil, err
}

func (db *DataStoreMongo) UpdateDevicesGroup(
	ctx context.Context,
	devIDs []model.DeviceID,
	group model.GroupName,
) (*model.UpdateResult, error) {
	database := db.client.Database(DbName)
	collDevs := database.Collection(DbDevicesColl)

	var filter = bson.M{DBKeyTenantID: tenantIDFromContext(ctx)}
	switch len(devIDs) {
	case 0:
		return &model.UpdateResult{}, nil
	case 1:
		filter[DbDevId] = devIDs[0]
	default:
		filter[DbDevId] = bson.M{"$in": devIDs}
	}
	update := bson.M{
		"$set": bson.M{
			model.AttrNameGroup: group,
		},
	}
	res, err := collDevs.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &model.UpdateResult{
		MatchedCount: res.MatchedCount,
		UpdatedCount: res.ModifiedCount,
	}, nil
}

// UpdateDeviceText updates the device text field
func (db *DataStoreMongo) UpdateDeviceText(
	ctx context.Context,
	deviceID model.DeviceID,
	text string,
) error {
	filter := bson.M{
		DbDevId:       deviceID.String(),
		DBKeyTenantID: tenantIDFromContext(ctx),
	}

	update := bson.M{
		"$set": bson.M{
			DbDevAttributesText: text,
		},
	}

	database := db.client.Database(DbName)
	collDevs := database.Collection(DbDevicesColl)

	_, err := collDevs.UpdateOne(ctx, filter, update)
	return err
}

func (db *DataStoreMongo) GetFiltersAttributes(
	ctx context.Context,
) ([]model.FilterAttribute, error) {
	database := db.client.Database(DbName)
	collDevs := database.Collection(DbDevicesColl)

	const DbCount = "count"

	cur, err := collDevs.Aggregate(ctx, []bson.M{
		// Sample up to 5,000 devices to get a representative sample
		{
			"$match": bson.M{DBKeyTenantID: tenantIDFromContext(ctx)},
		},
		{
			"$limit": FiltersAttributesMaxDevices,
		},
		{
			"$project": bson.M{
				"_id": 0,
				"attributes": bson.M{
					"$objectToArray": "$" + DbDevAttributes,
				},
			},
		},
		{
			"$unwind": "$" + DbDevAttributes,
		},
		{
			"$group": bson.M{
				DbDevId: bson.M{
					DbDevAttributesName:  "$" + DbDevAttributes + ".v." + DbDevAttributesName,
					DbDevAttributesScope: "$" + DbDevAttributes + ".v." + DbDevAttributesScope,
				},
				DbCount: bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$limit": FiltersAttributesLimit,
		},
		{
			"$sort": bson.D{
				{Key: DbCount, Value: -1},
				{Key: DbDevId + "." + DbDevAttributesScope, Value: 1},
				{Key: DbDevId + "." + DbDevAttributesName, Value: 1},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var attributes []model.FilterAttribute
	type Result struct {
		Group struct {
			Name  string `bson:"name"`
			Scope string `bson:"scope"`
		} `bson:"_id"`
		Count int32 `bson:"count"`
	}
	for cur.Next(ctx) {
		var elem Result
		err = cur.Decode(&elem)
		if err != nil {
			break
		}
		attributes = append(attributes, model.FilterAttribute{
			Name:  elem.Group.Name,
			Scope: elem.Group.Scope,
			Count: elem.Count,
		})
	}

	return attributes, nil
}

func (db *DataStoreMongo) DeleteGroup(
	ctx context.Context,
	group model.GroupName,
) (chan model.DeviceID, error) {
	deviceIDs := make(chan model.DeviceID)

	database := db.client.Database(DbName)
	collDevs := database.Collection(DbDevicesColl)

	filter := bson.M{
		model.AttrNameGroup: group,
		DBKeyTenantID:       tenantIDFromContext(ctx),
	}

	const batchMaxSize = 100
	batchSize := int32(batchMaxSize)
	findOptions := &mopts.FindOptions{
		Projection: bson.M{DbDevId: 1},
		BatchSize:  &batchSize,
	}
	cursor, err := collDevs.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	go func() {
		defer cursor.Close(ctx)
		batch := make([]model.DeviceID, batchMaxSize)
		batchSize := 0

		update := bson.M{"$unset": bson.M{model.AttrNameGroup: 1}}
		device := &model.Device{}
		defer close(deviceIDs)

	next:
		for {
			hasNext := cursor.Next(ctx)
			if !hasNext {
				if batchSize > 0 {
					break
				}
				return
			}
			if err = cursor.Decode(&device); err == nil {
				batch[batchSize] = device.ID
				batchSize++
				if len(batch) == batchSize {
					break
				}
			}
		}

		_, _ = collDevs.UpdateMany(ctx, bson.M{
			DbDevId: bson.M{"$in": batch[:batchSize],
				DBKeyTenantID: tenantIDFromContext(ctx),
			}}, update)
		for _, item := range batch[:batchSize] {
			deviceIDs <- item
		}
		batchSize = 0
		goto next
	}()

	return deviceIDs, nil
}

func (db *DataStoreMongo) UnsetDevicesGroup(
	ctx context.Context,
	deviceIDs []model.DeviceID,
	group model.GroupName,
) (*model.UpdateResult, error) {
	database := db.client.Database(DbName)
	collDevs := database.Collection(DbDevicesColl)

	var filter bson.D
	// Add filter on device id (either $in or direct indexing)
	switch len(deviceIDs) {
	case 0:
		return &model.UpdateResult{}, nil
	case 1:
		filter = bson.D{{Key: DbDevId, Value: deviceIDs[0]}}
	default:
		filter = bson.D{{Key: DbDevId, Value: bson.M{"$in": deviceIDs}}}
	}
	// Append filter on group
	filter = append(
		filter,
		bson.E{Key: model.AttrNameGroup, Value: group},
		bson.E{Key: DBKeyTenantID, Value: tenantIDFromContext(ctx)},
	)
	// Create unset operation on group attribute
	update := bson.M{
		"$unset": bson.M{
			model.AttrNameGroup: "",
		},
	}
	res, err := collDevs.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &model.UpdateResult{
		MatchedCount: res.MatchedCount,
		UpdatedCount: res.ModifiedCount,
	}, nil
}

func predicateToQuery(pred model.FilterPredicate) (bson.D, error) {
	if err := pred.Validate(); err != nil {
		return nil, err
	}
	name := fmt.Sprintf(
		"%s.%s-%s.value",
		DbDevAttributes,
		pred.Scope,
		model.GetDeviceAttributeNameReplacer().Replace(pred.Attribute),
	)
	return bson.D{{
		Key: name, Value: bson.D{{Key: pred.Type, Value: pred.Value}},
	}}, nil
}

func (db *DataStoreMongo) ListGroups(
	ctx context.Context,
	filters []model.FilterPredicate,
) ([]model.GroupName, error) {
	c := db.client.
		Database(DbName).
		Collection(DbDevicesColl)

	fltr := bson.D{
		{Key: model.AttrNameGroup, Value: bson.M{"$exists": true}},
		{Key: DBKeyTenantID, Value: tenantIDFromContext(ctx)},
	}
	if len(fltr) > 0 {
		for _, p := range filters {
			q, err := predicateToQuery(p)
			if err != nil {
				return nil, errors.Wrap(
					err, "store: bad filter predicate",
				)
			}
			fltr = append(fltr, q...)
		}
	}
	results, err := c.Distinct(
		ctx, model.AttrNameGroup, fltr,
	)
	if err != nil {
		return nil, err
	}

	groups := make([]model.GroupName, len(results))
	for i, d := range results {
		groups[i] = model.GroupName(d.(string))
	}
	return groups, nil
}

func (db *DataStoreMongo) GetDevicesByGroup(
	ctx context.Context,
	group model.GroupName,
	skip,
	limit int,
) ([]model.DeviceID, int, error) {
	c := db.client.
		Database(DbName).
		Collection(DbDevicesColl)

	filter := bson.M{
		model.AttrNameGroup: group,
		DBKeyTenantID:       tenantIDFromContext(ctx),
	}
	result := c.FindOne(ctx, filter)
	if result == nil {
		return nil, -1, store.ErrGroupNotFound
	}

	var dev model.Device
	err := result.Decode(&dev)
	if err != nil {
		return nil, -1, store.ErrGroupNotFound
	}

	hasGroup := group != ""
	devices, totalDevices, e := db.GetDevices(ctx,
		store.ListQuery{
			Skip:      skip,
			Limit:     limit,
			Filters:   nil,
			Sort:      nil,
			HasGroup:  &hasGroup,
			GroupName: string(group)})
	if e != nil {
		return nil, -1, errors.Wrap(e, "failed to get device list for group")
	}

	resIds := make([]model.DeviceID, len(devices))
	for i, d := range devices {
		resIds[i] = d.ID
	}
	return resIds, totalDevices, nil
}

func (db *DataStoreMongo) GetDeviceGroup(
	ctx context.Context,
	id model.DeviceID,
) (model.GroupName, error) {
	dev, err := db.GetDevice(ctx, id)
	if err != nil || dev == nil {
		return "", store.ErrDevNotFound
	}

	return dev.Group, nil
}

func (db *DataStoreMongo) DeleteDevices(
	ctx context.Context, ids []model.DeviceID,
) (*model.UpdateResult, error) {
	var filter = bson.M{
		DBKeyTenantID: tenantIDFromContext(ctx),
	}
	database := db.client.Database(DbName)
	collDevs := database.Collection(DbDevicesColl)

	switch len(ids) {
	case 0:
		// This is a no-op, don't bother requesting mongo.
		return &model.UpdateResult{DeletedCount: 0}, nil
	case 1:
		filter[DbDevId] = ids[0]
	default:
		filter[DbDevId] = bson.M{"$in": ids}
	}
	res, err := collDevs.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &model.UpdateResult{
		DeletedCount: res.DeletedCount,
	}, nil
}

func (db *DataStoreMongo) GetAllAttributeNames(ctx context.Context) ([]string, error) {
	c := db.client.Database(DbName).Collection(DbDevicesColl)

	match := bson.M{DBKeyTenantID: tenantIDFromContext(ctx)}
	result, err := c.Distinct(ctx, "inventory.name", match)
	if err != nil {
		return nil, err
	}
	attributeNames := make([]string, 0, len(result))
	for _, res := range result {
		if s, ok := res.(string); ok {
			attributeNames = append(attributeNames, s)
		}
	}
	return attributeNames, nil
}

func (db *DataStoreMongo) SearchDevices(
	ctx context.Context,
	searchParams model.SearchParams,
) ([]model.Device, int, error) {
	c := db.client.Database(DbName).Collection(DbDevicesColl)

	queryFilters := make([]bson.M, 0, len(searchParams.Filters)+1)
	queryFilters = append(queryFilters, bson.M{DBKeyTenantID: tenantIDFromContext(ctx)})
	mgoFilter := searchParams.Filters.ToMongoFilter()
	if mgoFilter != nil {
		queryFilters = append(queryFilters, mgoFilter)
	}

	// FIXME: remove after migrating ids to attributes
	if len(searchParams.DeviceIDs) > 0 {
		queryFilters = append(queryFilters, bson.M{"_id": bson.M{"$in": searchParams.DeviceIDs}})
	}

	if searchParams.Text != "" {
		queryFilters = append(queryFilters, bson.M{
			"$text": bson.M{
				"$search": utils.TextToKeywords(searchParams.Text),
			},
		})
	}

	findQuery := bson.M{"$and": queryFilters}

	findOptions := mopts.Find()
	findOptions.SetSkip(int64((searchParams.Page - 1) * searchParams.PerPage))
	findOptions.SetLimit(int64(searchParams.PerPage))

	if len(searchParams.Attributes) > 0 {
		var scopeProjections = map[string][]string{}
		for _, attribute := range searchParams.Attributes {
			scopeProjections[string(attribute.Scope)] = append(
				scopeProjections[string(attribute.Scope)],
				attribute.Attribute,
			)
		}
		projection := make(map[string]any, len(scopeProjections))
		for scope, in := range scopeProjections {
			projection[scope] = bson.M{
				"$filter": bson.M{
					"input": "$" + scope,
					"cond": bson.M{
						"$in": bson.A{"$$this.name", in},
					},
				},
			}
		}
		findOptions.SetProjection(projection)
	}

	if searchParams.Text != "" {
		findOptions.SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})
	} else if len(searchParams.Sort) > 0 {
		sortField := make(bson.D, len(searchParams.Sort))
		for i, sortQ := range searchParams.Sort {
			var field string
			if sortQ.Scope == model.AttrScopeIdentity && sortQ.Attribute == model.AttrNameID {
				field = DbDevId
			} else {
				name := fmt.Sprintf(
					"%s-%s",
					sortQ.Scope,
					model.GetDeviceAttributeNameReplacer().Replace(sortQ.Attribute),
				)
				field = fmt.Sprintf("%s.%s.value", DbDevAttributes, name)
			}
			sortField[i] = bson.E{Key: field, Value: 1}
			if sortQ.Order == "desc" {
				sortField[i].Value = -1
			}
		}
		findOptions.SetSort(sortField)
	}

	cursor, err := c.Find(ctx, findQuery, findOptions)
	if err != nil {
		return nil, -1, errors.Wrap(err, "failed to search devices")
	}
	defer cursor.Close(ctx)

	devices := []model.Device{}

	if err = cursor.All(ctx, &devices); err != nil {
		return nil, -1, errors.Wrap(err, "failed to search devices")
	}

	count, err := c.CountDocuments(ctx, findQuery)
	if err != nil {
		return nil, -1, errors.Wrap(err, "failed to search devices")
	}

	return devices, int(count), nil
}

func (db *DataStoreMongo) GetDeviceTierStatisticsByStatus(
	ctx context.Context,
) (
	*model.DeviceStatisticsByStatus,
	error,
) {
	collection := db.client.Database(DbName).Collection(DbDevicesColl)

	match := bson.M{
		"$match": bson.M{
			DBKeyTenantID: tenantIDFromContext(ctx),
			"status": bson.M{
				"$in": []string{
					model.DeviceStatusAccepted,
					model.DeviceStatusPending,
				},
			},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"status": "$status",
				"tier": bson.M{
					"$ifNull": bson.A{
						"$tier",
						tiers.StandardTier,
					},
				},
			},
			"devices": bson.M{
				"$sum": 1,
			},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"_id":     0,
			"status":  "$_id.status",
			"tier":    "$_id.tier",
			"devices": 1,
		},
	}

	cursor, err := collection.Aggregate(
		ctx,
		[]bson.M{match, group, project},
	)
	if err != nil {
		return nil, err
	}

	var deviceCounts []struct {
		Status  string `bson:"status"`
		Tier    string `bson:"tier"`
		Devices uint64 `bson:"devices"`
	}

	err = cursor.All(ctx, &deviceCounts)
	if err != nil {
		return nil, err
	}

	var statistics model.DeviceStatisticsByStatus
	for _, count := range deviceCounts {
		var dest *model.DeviceCountPerTier
		if count.Status == model.DeviceStatusAccepted {
			dest = &statistics.Accepted
		} else {
			dest = &statistics.Pending
		}

		switch count.Tier {
		case tiers.StandardTier:
			dest.Standard = count.Devices
		case tiers.MicroTier:
			dest.Micro = count.Devices
		case tiers.SystemTier:
			dest.System = count.Devices
		}
	}

	return &statistics, nil
}

func indexAttr(s *mongo.Client, ctx context.Context, attr string) error {
	l := log.FromContext(ctx)
	c := s.Database(DbName).Collection(DbDevicesColl)

	indexView := c.Indexes()
	keys := bson.D{
		{Key: "status", Value: 1},
		{Key: "$elemMatch", Value: bson.M{"name": attr}},
	}
	_, err := indexView.CreateOne(ctx, mongo.IndexModel{Keys: keys, Options: &mopts.IndexOptions{
		Name: &attr,
	}})

	if err != nil {
		if isTooManyIndexes(err) {
			l.Warnf(
				"failed to index attr %s in db %s: too many indexes",
				attr,
				DbName,
			)
		} else {
			return errors.Wrapf(
				err,
				"failed to index attr %s in db %s",
				attr,
				DbName,
			)
		}
	}

	return nil
}

func isTooManyIndexes(e error) bool {
	return strings.HasPrefix(e.Error(), "add index fails, too many indexes for inventory.devices")
}
