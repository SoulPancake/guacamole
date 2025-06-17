package main

/*
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <sys/types.h>

int pipefd[2]; // pipe file descriptors: [0] = read, [1] = write

int setup_pipe_and_fork() {
    if (pipe(pipefd) == -1) {
        perror("pipe");
        return -1;
    }

    pid_t pid = fork();
    if (pid == -1) {
        perror("fork");
        return -1;
    }

    if (pid == 0) {
        // Child process
        close(pipefd[0]); // Close read end
        printf("Child process writing to pipe...\n");
        const char *msg = "hello from child\n";
        write(pipefd[1], msg, strlen(msg));
        close(pipefd[1]);
        _exit(0);
    } else {
        // Parent process
        close(pipefd[1]); // Close write end
        return pipefd[0]; // return read end
    }
}
*/
import "C"
import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rfd := C.setup_pipe_and_fork()
	if rfd < 0 {
		fmt.Println("error setting up pipe and fork")
		return
	}

	go func() {
		fmt.Println("goroutine running while waiting for pipe data...")
	}()

	file := os.NewFile(uintptr(rfd), "pipe")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("Go received:", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading from pipe:", err)
	}
}
