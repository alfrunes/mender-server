// Copyright 2022 Northern.tech AS
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

	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-server/pkg/log"
	"github.com/mendersoftware/mender-server/pkg/mongo/migrate"

	"github.com/mendersoftware/mender-server/services/inventory/store"
)

// WithAutomigrate enables automatic migration and returns a new datastore based
// on current one
func (db *DataStoreMongo) WithAutomigrate() store.DataStore {
	return &DataStoreMongo{
		client:      db.client,
		automigrate: true,
	}
}

func (db *DataStoreMongo) Migrate(
	ctx context.Context,
	version string,
) error {
	l := log.FromContext(ctx)

	l.Infof("migrating %s", DbName)

	m := migrate.SimpleMigrator{
		Client:      db.client,
		Db:          DbName,
		Automigrate: db.automigrate,
	}

	ver, err := migrate.NewVersion(version)
	if err != nil {
		return errors.Wrap(err, "failed to parse service version")
	}

	migrations := []migrate.Migration{
		&migration_0_2_0{
			ms:  db,
			ctx: ctx,
		},
		&migration_1_0_0{
			ms:  db,
			ctx: ctx,
		},
		&migration_1_0_1{
			ms:  db,
			ctx: ctx,
		},
		&migration_1_0_2{
			ms:  db,
			ctx: ctx,
		},
		&migration_1_1_0{
			ms:  db,
			ctx: ctx,
		},
	}

	err = m.Apply(ctx, *ver, migrations)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}
