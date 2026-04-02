package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Service struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Version string `json:"version"`
	URL string `json:"url"`
	Type string `json:"type"`
	Status string `json:"status"`
	Tags string `json:"tags"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"registry.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS services(id TEXT PRIMARY KEY,name TEXT NOT NULL,version TEXT DEFAULT '',url TEXT DEFAULT '',type TEXT DEFAULT '',status TEXT DEFAULT 'active',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Service)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO services(id,name,version,url,type,status,tags,created_at)VALUES(?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Version,e.URL,e.Type,e.Status,e.Tags,e.CreatedAt);return err}
func(d *DB)Get(id string)*Service{var e Service;if d.db.QueryRow(`SELECT id,name,version,url,type,status,tags,created_at FROM services WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Version,&e.URL,&e.Type,&e.Status,&e.Tags,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Service{rows,_:=d.db.Query(`SELECT id,name,version,url,type,status,tags,created_at FROM services ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Service;for rows.Next(){var e Service;rows.Scan(&e.ID,&e.Name,&e.Version,&e.URL,&e.Type,&e.Status,&e.Tags,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM services WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM services`).Scan(&n);return n}
