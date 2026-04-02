package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-registry/internal/server";"github.com/stockyard-dev/stockyard-registry/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./registry-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("registry: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Registry — Self-hosted container and package registry\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("registry: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
