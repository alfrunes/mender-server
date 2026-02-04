package mongo

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"

	"github.com/mendersoftware/mender-server/services/inventory/model"
)

type Device struct {
	Identity     model.DeviceAttributes `bson:"identity"`
	Inventory    model.DeviceAttributes `bson:"inventory"`
	Tags         model.DeviceAttributes `bson:"tags"`
	model.Device `bson:",inline"`
}

var (
	tDeviceP = reflect.TypeOf(&model.Device{})
	tDevice  = reflect.TypeOf(model.Device{})
)

func deviceDecodeValue(ec bsoncodec.DecodeContext, r bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tDevice {
		return bsoncodec.ValueDecoderError{
			Name:     "DeviceDecodeValue",
			Types:    []reflect.Type{tDevice},
			Received: val,
		}
	}
	var (
		decodeDevice Device
		refDevice    = reflect.ValueOf(&decodeDevice).Elem()
	)
	dec, err := ec.LookupDecoder(refDevice.Type())
	if err != nil {
		return err
	}
	err = dec.DecodeValue(ec, r, refDevice)
	if err != nil {
		return err
	}
	for i := range decodeDevice.Identity {
		decodeDevice.Identity[i].Scope = model.AttrScopeIdentity
	}
	for i := range decodeDevice.Inventory {
		decodeDevice.Inventory[i].Scope = model.AttrScopeInventory
	}
	for i := range decodeDevice.Tags {
		decodeDevice.Tags[i].Scope = model.AttrScopeTags
	}
	decodeDevice.Attributes = decodeDevice.Identity
	if decodeDevice.Status != "" {
		decodeDevice.Attributes = append(decodeDevice.Attributes, model.DeviceAttribute{
			Name:  "status",
			Scope: model.AttrScopeIdentity,
			Value: decodeDevice.Status,
		})
	}
	decodeDevice.Attributes = append(decodeDevice.Attributes, decodeDevice.Inventory...)
	decodeDevice.Attributes = append(decodeDevice.Attributes, decodeDevice.Tags...)
	val.Set(reflect.ValueOf(decodeDevice.Device))
	return nil
}
