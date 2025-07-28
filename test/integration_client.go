package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	proto "zip-archive/protos"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect: ", err)
	}
	defer conn.Close()
	client := proto.NewArchiveZipServiceClient(conn)

	log.Println("Create Task:")
	createResp, err := client.CreateTask(context.Background(), &proto.CreateTaskRequest{})
	if err != nil {
		log.Fatal("createTask failed: ", err)
	}
	log.Printf("task created: id=%s, status=%s", createResp.TaskId, createResp.Status)
	taskID := createResp.TaskId

	log.Println("GetTaskStatus: ")
	statusResp, err := client.GetTaskStatus(context.Background(), &proto.GetTaskStatusRequest{
		TaskId: taskID,
	})
	if err != nil {
		log.Fatal("getTaskStatus failed: ", err)
	}
	log.Printf("task status: id=%s, status=%s", statusResp.TaskId, statusResp.Status)

	log.Println("AddFiles: ")
	urls := []string{
		"https://images.pexels.com/photos/32715939/pexels-photo-32715939.jpeg?cs=srgb&dl=pexels-willianjusten-32715939.jpg&fm=jpg",
		//"https://file-examples.com/wp-content/storage/2017/10/file-example_PDF_1MB.pdf",
		"https://images.pexels.com/photos/20514819/pexels-photo-20514819.jpeg?cs=srgb&dl=pexels-diana-reyes-227887231-20514819.jpg&fm=jpg",
	}

	addResp, err := client.AddFiles(context.Background(), &proto.AddFilesRequest{
		TaskId: taskID,
		Urls:   urls,
	})
	if err != nil {
		log.Fatal("addFiles failed: ", err)
	}
	log.Printf("files added: status=%s", addResp.Status)
	if len(addResp.FailedUrls) > 0 {
		log.Printf("failed URLs: %v", addResp.FailedUrls)
	}

	log.Println("GetTaskStatus after adding files")
	statusResp, err = client.GetTaskStatus(context.Background(), &proto.GetTaskStatusRequest{
		TaskId: taskID,
	})
	if err != nil {
		log.Fatal("getTaskStatus failed: ", err)
	}
	log.Printf("task status: id=%s, status=%s", statusResp.TaskId, statusResp.Status)
	if len(statusResp.FailedUrls) > 0 {
		log.Printf("failed URLs: %v", statusResp.FailedUrls)
	}

	log.Println("Add the third file")
	finalUrls := []string{
		"https://www.pexels.com/ru-ru/photo/33054756/",
	}
	addResp, err = client.AddFiles(context.Background(), &proto.AddFilesRequest{
		TaskId: taskID,
		Urls:   finalUrls,
	})
	if err != nil {
		log.Fatal("addFiles to finalUrls failed: ", err)
	}
	log.Printf("file added: status=%s", addResp.Status)

	log.Println("Archive's processing")
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)

		statusResp, err = client.GetTaskStatus(context.Background(), &proto.GetTaskStatusRequest{
			TaskId: taskID,
		})
		if err != nil {
			log.Fatal("getTaskStatus failed: ", err)
		}
		log.Printf("status check %d: %s", i+1, statusResp.Status)

		if statusResp.Status == "completed" || statusResp.Status == "failed" {
			break
		}
	}

	if statusResp.Status == "completed" {
		log.Println("DownloadArchive: ")
		downloadResp, err := client.DownloadArchive(context.Background(), &proto.DownloadArchiveRequest{
			TaskId: taskID,
		})
		if err != nil {
			log.Fatal("downloadArchive failed: ", err)
		}
		log.Printf("archive downloaded: filename=%s, size=%d bytes",
			downloadResp.Filename, len(downloadResp.ArchiveData))
	} else {
		log.Printf("task isn't completed. status: %s", statusResp.Status)
	}

	_, err = client.GetTaskStatus(context.Background(), &proto.GetTaskStatusRequest{
		TaskId: "fsadjksdfkj",
	})
	if err != nil {
		log.Printf("expected error: %v", err)
	} else {
		log.Println("smth happened")
	}

	log.Println("ALL TEST COMPLETED!!!")

}
