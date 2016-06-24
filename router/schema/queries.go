package schema

import (
	"github.com/flynn/que-go"
	"github.com/jackc/pgx"
)

var preparedStatements = map[string]string{
	"ping":             ping,
	"insert_tcp_route": insertTcpRoute,
}

func PrepareStatements(conn *pgx.Conn) error {
	for name, sql := range preparedStatements {
		if _, err := conn.Prepare(name, sql); err != nil {
			return err
		}
	}
	if err := que.PrepareStatements(conn); err != nil {
		return err
	}
	return nil
}

const (
	// misc
	ping = `SELECT 1`

	// tcp
	insertTcpRoute = `
	INSERT INTO tcp_routes (parent_ref, service, leader, port)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`

	// http
	insertHttpRoute = `
	INSERT INTO http_routes (parent_ref, service, leader, domain, sticky, path)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`

	// certificates
	selectCertBySha = `
	SELECT id, created_at, updated_at FROM certificates
	WHERE cert_sha256 = $1 AND deleted_at IS NULL`

	selectCertByRouteId = `
	SELECT c.id, c.cert, c.key, c.created_at, c.updated_at, ARRAY(
		SELECT http_route_id::varchar FROM route_certificates
		WHERE certificate_id = $1
	) FROM certificates AS c WHERE c.id = $1`

	listCerts = `
	SELECT c.id, c.cert, c.key, c.created_at, c.updated_at, ARRAY(
		SELECT http_route_id::varchar FROM route_certificates AS rc
		WHERE rc.certificate_id = c.id
	) FROM certificates AS c`

	listCertRoutes = `
	SELECT r.id, r.parent_ref, r.service, r.leader, r.domain, r.sticky, r.path, r.created_at, r.updated_at FROM http_routes AS r
	INNER JOIN route_certificates AS rc ON rc.http_route_id = r.id AND rc.certificate_id = $1`

	insertCert = `
	INSERT INTO certificates (cert, key, cert_sha256)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	insertRouteCertificate = `
	INSERT INTO route_certificates (http_route_id, certificate_id)
	VALUES ($1, $2)
	`

	deleteRouteCertificateByCertifcateId = `
	DELETE FROM route_certificates
	WHERE certificate_id = $1`

	deleteRouteCertificateByRouteId = `
	DELETE FROM route_certificates
	WHERE http_route_id = $1`
)
