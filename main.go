package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
	u "github.com/hughluo/go-tiny-url-kgs/utils"
	"github.com/hughluo/go-tiny-url/pb"
	UTILS "github.com/hughluo/go-tiny-url/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCServer struct{}

var CLIENT *redis.Client

func main() {
	UTILS.ConfigureLog()
	// Configure redis client
	CLIENT = createClient()

	if INIT_REDIS_FREE := UTILS.GetEnv("INIT_REDIS_FREE", "true"); INIT_REDIS_FREE == "true" {
		fmt.Println("INIT_REDIS_FREE is true or undefined, start initRedis")
		initRedisFree()
		fmt.Printf("REDIS_FREE inited, free tinyURL amount: %d\n", getSetFreeAmount())
	}

	// Set up gRPC
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		fmt.Printf("Failed to listen:  %v", err)
		log.Fatalf("Failed to listen:  %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterKGSServiceServer(s, &gRPCServer{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve:  %v", err)
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *gRPCServer) GetFreeGoTinyURL(cxt context.Context, req *pb.KGSRequest) (*pb.KGSResponse, error) {
	result := &pb.KGSResponse{}
	result.Result = popSetFree()
	fmt.Printf("KGS req: %s result: %s\n", req, result.Result)
	logMessage := fmt.Sprintf("KGS req: %s result: %s", req, result.Result)
	log.Println(logMessage)
	return result, nil
}

func initRedisFree() {
	KEY_LENGTH, err := strconv.Atoi(UTILS.GetEnv("KEY_LENGTH", "4"))
	if err != nil {
		KEY_LENGTH = 4
		fmt.Println("KEY_LENGTH not valid, fallback to 4")
	}
	base62 := u.GetBase62String()
	base62Slice := strings.Split(base62, "")
	addAllTinyURLToSetFree(base62Slice, KEY_LENGTH)
}

func addAllTinyURLToSetFree(charArray []string, keyLength int) {
	addAllTinyURLToSetFreeHelper(charArray, len(charArray), "", keyLength)
}

func addAllTinyURLToSetFreeHelper(charArray []string, n int, prefix string, length int) {
	if length == 0 {
		addToSetFree(prefix)
		return
	}

	for index := 0; index < n; index++ {
		newPrefix := prefix + charArray[index]
		addAllTinyURLToSetFreeHelper(charArray, n, newPrefix, length-1)
	}
}

func addToSetFree(freeTinyURL string) {
	err := CLIENT.SAdd("FREE", freeTinyURL).Err()
	if err != nil {
		panic(err)
	}
}

func getSetFreeAmount() int64 {
	amount, err := CLIENT.SCard("FREE").Result()
	if err != nil {
		panic(err)
	}
	return amount
}

func popSetFree() string {
	freeTinyURL, err := CLIENT.SPop("FREE").Result()
	if err != nil {
		panic(err)
	}
	return freeTinyURL
}

func createClient() *redis.Client {
	REDIS_FREE_PASSWORD := UTILS.GetEnv("REDIS_FREE_PASSWORD", "")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-free-service:6379",
		Password: REDIS_FREE_PASSWORD,
		DB:       0, // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Fail to connect to redis")
		log.Fatalf("Fail to connect to redis")
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}
