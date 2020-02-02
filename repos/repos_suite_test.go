package repos_test

import (
	"database/sql"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tv2169145/golang-grpc/repos"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	err error
	db *xorm.Engine
	dbSql *sql.DB
	mock sqlmock.Sqlmock

	gr repos.GlobalRepository

	truncateUsers = func() {
		mock.ExpectQuery("TRUNCATE users").WillReturnRows(sqlmock.NewRows([]string{}))
		_, err = db.Query("TRUNCATE users")
		Ω(err).To(BeNil())
	}

	clearDatabase = func() {
		if db == nil {
			Fail("unable to run test because database is missing")
		}
		truncateUsers()
		return
	}
)

var (
	_ = BeforeSuite(func() {
		// connection string - root:pass@tcp(localhost:3306)/grpc
		// root:12345678@tcp(localhost:3306)/grpc
		db, err = xorm.NewEngine("mysql", "")
		Ω(err).To(BeNil())
		dbSql, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		Ω(err).To(BeNil())
		db.DB().DB = dbSql
		gr = repos.GlobalRepo(db)
	})

	_ = AfterSuite(func() {
		err = mock.ExpectationsWereMet()
		Ω(err).To(BeNil())
	})
)

func TestRepos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repos Suite")
}
