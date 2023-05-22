package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	// "time"

	"github.com/kardianos/service"
)

var logger service.Logger

type Config struct {
	Name, DisplayName, Description string

	Dir  string
	Exec string
	Args []string
	Env  []string

	Stderr, Stdout string
}

type program struct {
	exit    chan struct{}
	service service.Service

	*Config

	cmd *exec.Cmd
}

// func (p *program) runMindoc() {
// 	// 获取命令行参数
// 	cwd := os.Args[1] // mindoc运行目录
// 	cmd := os.Args[2] // mindoc命令

// 	// 连续尝试次数
// 	retries := 0
// 	maxRetries := 5

// 	for {
// 		// 执行mindoc命令
// 		cmd := exec.Command(cmd, cwd)
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		err := cmd.Start()
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		// cmd.Process.Kill()

// 		// 等待mindoc运行结束,或发生错误
// 		err = cmd.Wait()
// 		if err != nil {
// 			retries++
// 			fmt.Printf("mindoc exited with error, retrying (%d/%d)...\n", retries, maxRetries)
// 			time.Sleep(5 * time.Second)
// 		} else {
// 			retries = 0
// 		}

// 		// 超过最大重试次数,退出
// 		if retries > maxRetries {
// 			fmt.Println(" mindoc max retries exceeded, exiting!")
// 			return
// 		}
// 	}
// }

func (p *program) run() {
	logger.Info("Starting ", p.DisplayName)
	defer func() {
		if service.Interactive() {
			p.Stop(p.service)
		} else {
			p.service.Stop()
		}
	}()

	if p.Stderr != "" {
		f, err := os.OpenFile(p.Stderr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			logger.Warningf("Failed to open std err %q: %v", p.Stderr, err)
			return
		}
		defer f.Close()
		p.cmd.Stderr = f
	} else {
		p.cmd.Stderr = os.Stderr
	}
	if p.Stdout != "" {
		f, err := os.OpenFile(p.Stdout, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			logger.Warningf("Failed to open std out %q: %v", p.Stdout, err)
			return
		}
		defer f.Close()
		p.cmd.Stdout = f
	} else {
		p.cmd.Stdout = os.Stdout
	}

	err := p.cmd.Run()
	if err != nil {
		logger.Warningf("Error running: %v", err)
	}

	return
}

func (p *program) Start(s service.Service) error {
	// Look for exec.
	// Verify home directory.
	fullExec, err := exec.LookPath(p.Exec)
	if err != nil {
		return fmt.Errorf("Failed to find executable %q: %v", p.Exec, err)
	}

	p.cmd = exec.Command(fullExec, p.Args...)
	p.cmd.Dir = p.Dir
	p.cmd.Env = append(os.Environ(), p.Env...)
	// 启动mindoc-daemon
	go p.run()
	return nil
}

// func (p *program) Restart(s service.Service) error {
// 	return nil
// }

func (p *program) Stop(s service.Service) error {
	// 停止mindoc-daemon
	close(p.exit)
	logger.Info("Stopping ", p.DisplayName)
	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

// func (p *program) Install(s service.Service) error {
// 	return nil
// }

// func (p *program) Uninstall(s service.Service) error {
// 	return nil
// }

func getConfigPath() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, execname := filepath.Split(fullexecpath)
	ext := filepath.Ext(execname)
	name := execname[:len(execname)-len(ext)]

	return filepath.Join(dir, name+".json"), nil
}

func getConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &Config{}

	r := json.NewDecoder(f)
	err = r.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func main() {
	log.Println("mindoc daemon")
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Executable() failed: %v\n", err)
		return
	}
	log.Printf("executable: %s\n", executable)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Getwd() failed: %v\n", err)
		return
	}
	log.Printf("dir: %s\n", dir)

	/*
		svcConfig := &service.Config{
			Name:        "mindoc-daemon",
			DisplayName: "mindoc Daemon Service",
			Description: "Service to start and monitor mindoc process",
		}

		prg := &program{}
		s, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
		}
		if len(os.Args) > 1 {
			err = service.Control(s, os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		logger, err = s.Logger(nil)
		if err != nil {
			log.Fatal(err)
		}
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
	*/

	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	config, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	svcConfig := &service.Config{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		Description: config.Description,
	}

	prg := &program{
		exit: make(chan struct{}),

		Config: config,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.service = s

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
