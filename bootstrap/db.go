package bootstrap

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func (c *Container) GetDBRead() *pgx.Conn {
	if c.dBR != nil {
		c.dBR = c.getDB("read")
	}

	return c.dBR
}

func (c *Container) GetDBWrite() *pgx.Conn {
	if c.dBW != nil {
		c.dBW = c.getDB("write")
	}
	return c.dBW
}

func (c *Container) getDB(dbType string) *pgx.Conn {
	conn, err := pgx.Connect(c.ctx, viper.GetString("database."+dbType+".connstring"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return conn
}