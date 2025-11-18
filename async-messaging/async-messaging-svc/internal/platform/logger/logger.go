package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
* Ubicación del archivo de logs
 */
/* var ejecutable, _ = os.Executable()
var rutaAbsoluta = filepath.Join(filepath.Dir(ejecutable), "logs") */

var rutaAbsoluta = "/tmp/logs"
var fileName string

type ResponseRecorder struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Body.Write(b) // Capturar el contenido de la respuesta
	return r.ResponseWriter.Write(b)
}

// Crea el directorio en caso de que no exista
func setupLogsDirectory() {
	if _, err := os.Stat(rutaAbsoluta); os.IsNotExist(err) {
		err = os.MkdirAll(rutaAbsoluta, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func LoggerInfo() *logrus.Entry {
	setupLogsDirectory()

	if fileName == "" {
		fileName = "logger-" + time.Now().Format("02012006") + ".log"
	}

	file, err := os.OpenFile(rutaAbsoluta+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Println("Error al abrir el archivo de registro: ", err)
		return nil
	}

	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logEntry := logrus.WithFields(logrus.Fields{})

	return logEntry
}

func LoggerError() *logrus.Entry {
	setupLogsDirectory()

	if fileName == "" {
		fileName = "logger-" + time.Now().Format("02012006") + ".log"
	}

	file, err := os.OpenFile(rutaAbsoluta+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Println("Error al abrir el archivo de registro: ", err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyFile:  "file",
			logrus.FieldKeyFunc:  "func",
		},
	})
	logrus.SetReportCaller(true)

	// Crear un io.MultiWriter para dirigir la salida tanto al archivo como a la consola
	consoleWriter := os.Stdout
	multiWriter := io.MultiWriter(file, consoleWriter)
	logrus.SetOutput(multiWriter)

	//Setear parámetros del Mensaje
	logEntry := logrus.WithFields(logrus.Fields{})

	// escribe el error en el Parámetro msg
	return logEntry
}

func LoggerHttp(c *gin.Context, responseBody string) {
	setupLogsDirectory()
	fileName := "http_request-" + time.Now().Format("02012006") + ".log"
	file, err := os.OpenFile(rutaAbsoluta+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Println("Error al abrir el archivo de registro: ", err)
		return
	}

	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	//Setear parámetros del Mensaje
	logEntry := logrus.WithFields(logrus.Fields{
		"method":        c.Request.Method,
		"url":           c.Request.URL.String(),
		"headers":       c.Request.Header,
		"params":        c.Request.URL.Query(),
		"request_body":  string(bodyBytes),
		"response_body": responseBody,
		"status":        c.Writer.Status(),
	})

	if c.Writer.Status() != 200 {
		logEntry.Info("httpRequest")
	}
}
