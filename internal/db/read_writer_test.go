package db

import (
	"context"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/watchtower/internal/db/db_test"
	"github.com/hashicorp/watchtower/internal/oplog"
	"gotest.tools/assert"
)

func TestGormReadWriter_Update(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("simple", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)

		var foundUser db_test.TestUser
		foundUser.Id = user.Id
		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)

		user.FriendlyName = "friendly-" + id
		err = w.Update(context.Background(), user, []string{"FriendlyName"})
		assert.NilError(t, err)

		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.FriendlyName, foundUser.FriendlyName)
	})
	t.Run("valid-WithOplog", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)

		var foundUser db_test.TestUser
		foundUser.Id = user.Id
		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)

		user.FriendlyName = "friendly-" + id
		_, err = w.DoTx(
			context.Background(),
			20,
			ExpBackoff{},
			func(Writer) error {
				return w.Update(context.Background(), user, []string{"FriendlyName"},
					WithOplog(true),
					WithWrapper(InitTestWrapper(t)),
					WithMetadata(oplog.Metadata{
						"key-only":   nil,
						"deployment": []string{"amex"},
						"project":    []string{"central-info-systems", "local-info-systems"},
					}),
				)
			})
		assert.NilError(t, err)

		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.FriendlyName, foundUser.FriendlyName)
	})
	t.Run("nil-tx", func(t *testing.T) {
		w := GormReadWriter{Tx: nil}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Update(context.Background(), user, nil)
		assert.Check(t, err != nil)
		assert.Equal(t, err.Error(), "update tx is nil")
	})
	t.Run("no-wrapper-WithOplog", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)

		var foundUser db_test.TestUser
		foundUser.Id = user.Id
		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)

		user.FriendlyName = "friendly-" + id
		_, err = w.DoTx(
			context.Background(),
			20,
			ExpBackoff{},
			func(Writer) error {
				return w.Update(context.Background(), user, []string{"FriendlyName"},
					WithOplog(true),
					WithMetadata(oplog.Metadata{
						"key-only":   nil,
						"deployment": []string{"amex"},
						"project":    []string{"central-info-systems", "local-info-systems"},
					}),
				)
			})
		assert.Check(t, err != nil)
		assert.Equal(t, err.Error(), "error wrapper is nil for WithWrapper")
	})
	t.Run("no-metadata-WithOplog", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)

		var foundUser db_test.TestUser
		foundUser.Id = user.Id
		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)

		user.FriendlyName = "friendly-" + id
		_, err = w.DoTx(
			context.Background(),
			20,
			ExpBackoff{},
			func(Writer) error {
				return w.Update(context.Background(), user, []string{"FriendlyName"},
					WithOplog(true),
					WithWrapper(InitTestWrapper(t)),
				)
			})
		assert.Check(t, err != nil)
		assert.Equal(t, err.Error(), "error no metadata for WithOplog")
	})
}

func TestGormReadWriter_Create(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("simple", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)
		assert.Check(t, user.GetCreateTime() != nil)
		assert.Check(t, user.GetUpdateTime() != nil)

		var foundUser db_test.TestUser
		foundUser.Id = user.Id
		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)
	})
	t.Run("valid-WithOplog", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		_, err = w.DoTx(
			context.Background(),
			3,
			ExpBackoff{},
			func(Writer) error {
				return w.Create(
					context.Background(),
					user,
					WithOplog(true),
					WithWrapper(InitTestWrapper(t)),
					WithMetadata(oplog.Metadata{
						"key-only":   nil,
						"deployment": []string{"amex"},
						"project":    []string{"central-info-systems", "local-info-systems"},
					}),
				)
			})
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)

		var foundUser db_test.TestUser
		foundUser.Id = user.Id
		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)
	})
	t.Run("no-wrapper-WithOplog", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		_, err = w.DoTx(
			context.Background(),
			3,
			ExpBackoff{},
			func(Writer) error {
				return w.Create(
					context.Background(),
					user,
					WithOplog(true),
					WithMetadata(oplog.Metadata{
						"key-only":   nil,
						"deployment": []string{"amex"},
						"project":    []string{"central-info-systems", "local-info-systems"},
					}),
				)
			})
		assert.Check(t, err != nil)
		assert.Equal(t, err.Error(), "error wrapper is nil for WithWrapper")
	})
	t.Run("no-metadata-WithOplog", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		_, err = w.DoTx(
			context.Background(),
			3,
			ExpBackoff{},
			func(Writer) error {
				return w.Create(
					context.Background(),
					user,
					WithOplog(true),
					WithWrapper(InitTestWrapper(t)),
				)
			})
		assert.Check(t, err != nil)
		assert.Equal(t, err.Error(), "error no metadata for WithOplog")
	})
	t.Run("nil-tx", func(t *testing.T) {
		w := GormReadWriter{Tx: nil}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Create(context.Background(), user)
		assert.Check(t, err != nil)
		assert.Equal(t, err.Error(), "create tx is nil")
	})
}

func TestGormReadWriter_LookupByInternalId(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("simple", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)

		var foundUser db_test.TestUser
		foundUser.Id = user.Id
		err = w.LookupById(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)
	})
}

func TestGormReadWriter_LookupByFriendlyName(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("simple", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		user.FriendlyName = "fn-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.Id != 0)

		var foundUser db_test.TestUser
		foundUser.FriendlyName = "fn-" + id
		err = w.LookupByFriendlyName(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)
	})
}

func TestGormReadWriter_LookupByPublicId(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("simple", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		user.FriendlyName = "fn-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.PublicId != "")

		var foundUser db_test.TestUser
		foundUser.PublicId = user.PublicId
		err = w.LookupByPublicId(context.Background(), &foundUser)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)
	})
}

func TestGormReadWriter_LookupBy(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("simple", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		user.FriendlyName = "fn-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.PublicId != "")

		var foundUser db_test.TestUser
		err = w.LookupBy(context.Background(), &foundUser, "public_id = ?", user.PublicId)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUser.Id)
	})
}

func TestGormReadWriter_SearchBy(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("simple", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		id, err := uuid.GenerateUUID()
		assert.NilError(t, err)
		user, err := db_test.NewTestUser()
		assert.NilError(t, err)
		user.Name = "foo-" + id
		user.FriendlyName = "fn-" + id
		err = w.Create(context.Background(), user)
		assert.NilError(t, err)
		assert.Check(t, user.PublicId != "")

		var foundUsers []db_test.TestUser
		err = w.SearchBy(context.Background(), &foundUsers, "public_id = ?", user.PublicId)
		assert.NilError(t, err)
		assert.Equal(t, user.Id, foundUsers[0].Id)
	})
}

func TestGormReadWriter_Dialect(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("valid", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		d, err := w.Dialect()
		assert.NilError(t, err)
		assert.Equal(t, d, "postgres")
	})
	t.Run("nil-tx", func(t *testing.T) {
		w := GormReadWriter{Tx: nil}
		d, err := w.Dialect()
		assert.Check(t, err != nil)
		assert.Equal(t, d, "")
		assert.Equal(t, err.Error(), "create tx is nil for dialect")
	})
}

func TestGormReadWriter_DB(t *testing.T) {
	StartTest()
	t.Parallel()
	cleanup, url := SetupTest(t, "migrations/postgres")
	defer cleanup()
	defer CompleteTest() // must come after the "defer cleanup()"
	conn, err := TestConnection(url)
	assert.NilError(t, err)
	defer conn.Close()
	t.Run("valid", func(t *testing.T) {
		w := GormReadWriter{Tx: conn}
		d, err := w.DB()
		assert.NilError(t, err)
		assert.Check(t, d != nil)
		err = d.Ping()
		assert.NilError(t, err)
	})
	t.Run("nil-tx", func(t *testing.T) {
		w := GormReadWriter{Tx: nil}
		d, err := w.DB()
		assert.Check(t, err != nil)
		assert.Check(t, d == nil)
		assert.Equal(t, err.Error(), "create tx is nil for db")
	})
}
