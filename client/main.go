package main
 
import (
   "flag"
   "log"
   "net/http"
 
   pb "gRPC_GO_crud/proto"
   "github.com/gin-gonic/gin"
   "google.golang.org/grpc"
   "google.golang.org/grpc/credentials/insecure"
)
 
var (
   addr = flag.String("addr", "localhost:50051", "the address to connect to")
)
 
type User struct {
   ID    string `json:"id"`
   Name string `json:"Name"`
   Age int32 `json:"age"`
}
 
func main() {
   flag.Parse()
   conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
 
   if err != nil {
       log.Fatalf("did not connect: %v", err)
   }
 
   defer conn.Close()
   client := pb.NewUserServiceClient(conn)
 
   r := gin.Default()
   r.GET("/users/:id", func(ctx *gin.Context) {
       id := ctx.Param("id")
       res, err := client.GetUser(ctx, &pb.ReadUserRequest{Id: id})
       if err != nil {
           ctx.JSON(http.StatusNotFound, gin.H{
               "message": err.Error(),
           })
           return
       }
       ctx.JSON(http.StatusOK, gin.H{
           "user": res.User,
       })
   })
   r.POST("/users", func(ctx *gin.Context) {
       var user User
 
       err := ctx.ShouldBind(&user)
       if err != nil {
           ctx.JSON(http.StatusBadRequest, gin.H{
               "error": err,
           })
           return
       }
       data := &pb.User{
           Name: user.Name,
           Age: user.Age,
       }
       res, err := client.CreateUser(ctx, &pb.CreateUserRequest{
           User: data,
       })
       if err != nil {
           ctx.JSON(http.StatusBadRequest, gin.H{
               "error": err,
           })
           return
       }
       ctx.JSON(http.StatusCreated, gin.H{
           "user": res.User,
       })
   })
   r.PUT("/users/:id", func(ctx *gin.Context) {
	   id := ctx.Param("id")
       var user User
       err := ctx.ShouldBind(&user)
       if err != nil {
           ctx.JSON(http.StatusBadRequest, gin.H{
               "error": err.Error(),
           })
           return
       }
       res, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
           User: &pb.User{
               Id:    id,
               Name: user.Name,
           },
       })
       if err != nil {
           ctx.JSON(http.StatusBadRequest, gin.H{
               "error": err.Error(),
           })
           return
       }
       ctx.JSON(http.StatusOK, gin.H{
           "user": res.User,
       })
       return
 
   })
   r.Run(":5000")
 
}
