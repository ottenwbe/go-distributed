package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
)

var _ = Describe("Client", func() {

	var (
		mockStd *mockStdInOut
		client  *Client
		server  *Server
	)

	const (
		testAddress = ":8888"
	)

	Context("Client", func() {
		BeforeEach(func() {
			client = NewClient()
		})

		It("allows close to be called when not connected", func() {
			err := client.Close()
			Expect(err).To(BeNil())
		})
	})

	Context("Server", func() {
		BeforeEach(func() {
			server = NewServer(testAddress)
		})

		It("created a server object and set the test address", func() {
			Expect(server).To(Not(BeNil()))
			Expect(server.listenAddress).To(Equal(testAddress))
		})

		It("allows close to be called when not started", func() {
			err := server.Close()
			Expect(err).To(BeNil())
		})
	})

	Context("Client and Server Communication", func() {

		BeforeEach(func() {
			server = NewServer(testAddress)
			err := server.Start()
			panicOnError(err)

			client = NewClient()
			err = client.Connect(testAddress)
			panicOnError(err)
		})

		AfterEach(func() {
			err := client.Close()
			panicOnError(err)
			err = server.Close()
			panicOnError(err)
		})

		It("can send a message and receive it", func() {
			const expectedOutput = "Hi\n"

			actualOutput, err := client.CallServer("iH")
			failOnError(err)

			Expect(actualOutput).To(Equal(expectedOutput))
		})

		It("closes if server is not available", func() {

			_, err := client.CallServer("no msg")
			failOnError(err)

			err = server.Close()
			failOnError(err)

			_, err = client.CallServer("no msg")

			Expect(err).To(Not(BeNil()))
			Expect(err).To(Equal(io.EOF))
		})

	})

	Context("Utils", func() {
		It("can exit the app on error", func() {
			capturedCode := -1
			exitFunc = func(code int) {
				capturedCode = code
			}

			exitOnError(io.EOF, 1)

			Expect(capturedCode).To(Equal(1))
		})

		It("can reverse strings", func() {
			expectedText := "sgnirts"
			actualText := reverseString("strings")
			Expect(actualText).To(Equal(expectedText))
		})
	})

	Context("User Interaction", func() {
		BeforeEach(func() {
			// mock stdin and stdout
			mockStd = &mockStdInOut{}
			err := mockStd.mock()
			panicOnError(err)
		})

		AfterEach(func() {
			mockStd.Close()
		})

		It("can send messages to a server that have been given by a users", func() {
			const expectedOutput = "Client Server Demo App\n>> Type 'quit' to stop the client (or ctrl+D)\n>> Message to server: >> Message from server: iH\n>> Message to server: "

			err := mockStd.FinalWriteStdIn("Hi\nquit\n")
			panicOnError(err)

			main()

			actualOutput, err := mockStd.FinalReadAllStdOut()
			panicOnError(err)

			Expect(actualOutput).To(Equal(expectedOutput))
		})

	})

})

type mockStdInOut struct {
	mockStdin        *os.File
	mockStdout       *os.File
	mockStdinWriter  *os.File
	mockStdoutReader *os.File

	realStdout *os.File
	realStdin  *os.File
}

func (m *mockStdInOut) mock() error {
	var err error
	m.realStdout = os.Stdout
	m.realStdin = os.Stdin
	if m.mockStdoutReader, m.mockStdout, err = os.Pipe(); err != nil {
		return err
	}
	os.Stdout = m.mockStdout
	if m.mockStdin, m.mockStdinWriter, err = os.Pipe(); err != nil {
		return err
	}
	os.Stdin = m.mockStdin
	return nil
}

func (m *mockStdInOut) Close() {
	_ = m.CloseStdIn()
	_ = m.CloseStdOut()
	_ = m.mockStdoutReader.Close
	_ = m.mockStdin.Close()
	if m.realStdout != nil {
		os.Stdout = m.realStdout
	}
	if m.realStdin != nil {
		os.Stdin = m.realStdin
	}
}

func (m *mockStdInOut) CloseStdIn() error {
	if m.mockStdinWriter != nil {
		return m.mockStdinWriter.Close()
	}
	return errors.New("Could not close mocked stdin, it does not exist")
}

func (m *mockStdInOut) CloseStdOut() error {
	if m.mockStdout != nil {
		return m.mockStdout.Close()
	}
	return errors.New("Could not close mocked stdout, it does not exist")
}

func (m *mockStdInOut) FinalReadAllStdOut() (string, error) {
	var (
		err          error
		actualOutput []byte
	)
	// Read everything that goes to stdout
	if m.mockStdout != nil && m.mockStdoutReader != nil {
		err = m.CloseStdOut()
		if err == nil {
			actualOutput, err = ioutil.ReadAll(m.mockStdoutReader)
		}
	} else {
		err = errors.New("Could not read from non existent mocked stdout")
	}
	return string(actualOutput), err
}

func (m *mockStdInOut) FinalWriteStdIn(s string) (err error) {
	if m.mockStdinWriter != nil {
		_, err = m.mockStdinWriter.WriteString(s)
		if err == nil {
			err = m.mockStdinWriter.Close()
		}
	}
	return
}

func panicOnError(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func failOnError(e error) {
	if e != nil {
		Fail(e.Error())
	}
}
