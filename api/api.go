package api

import (
	"context"
	"fmt"
	converter "github.com/krol44/capy-converter/pkg"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"
)

func Run() {
	listener, err := net.Listen("tcp", ":3003")
	if err != nil {
		log.Fatalln(err)
	}

	mb, err := strconv.Atoi(os.Getenv("MAX_FILE_SIZE_MB"))
	if err != nil {
		log.Fatalln("error - MAX_FILE_SIZE_MB")
	}

	s := grpc.NewServer(grpc.MaxSendMsgSize(mb<<20), grpc.MaxRecvMsgSize(mb<<20),
		grpc.ConnectionTimeout(time.Minute*10))

	log.Info("Running...")

	converter.RegisterConverterServer(s, &Server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed: %v", err)
	}

}

type Server struct {
	converter.UnimplementedConverterServer
}

func (s *Server) GifToWebM(_ context.Context, object *converter.GifToWebMType) (*converter.GifToWebMType, error) {
	gifQni := UniqueId("file-convert-gif")
	gif, err := os.CreateTemp("/tmp", gifQni)
	if err != nil {
		return nil, err
	}
	defer func(path string) {
		err := os.Remove(path)
		if err != nil {
			log.Error(err)
		}
	}(gif.Name())

	_, err = gif.Write(object.File)
	if err != nil {
		return nil, err
	}

	webm := UniqueId("/tmp/file-convert") + ".webm"
	defer func(path string) {
		err := os.Remove(path)
		if err != nil {
			log.Error(err)
		}
	}(webm)

	_, err = gif.Write(object.File)
	if err != nil {
		return nil, err
	}

	err = exec.Command("ffmpeg",
		"-protocol_whitelist", "file",
		"-i", gif.Name(),
		"-c", "vp9", "-b:v", "0", "-crf", "30",
		"-y", webm).Run()

	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(webm)
	if err != nil {
		return nil, err
	}

	return &converter.GifToWebMType{File: file}, nil
}

func UniqueId(prefix string) string {
	now := time.Now()
	sec := now.Unix()
	use := now.UnixNano() % 0x100000
	return fmt.Sprintf("%s-%08x%05x", prefix, sec, use)
}

func LogSetup() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
	if l, err := log.ParseLevel("debug"); err == nil {
		log.SetLevel(l)
		log.SetReportCaller(l == log.DebugLevel)
		log.SetOutput(os.Stdout)
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
