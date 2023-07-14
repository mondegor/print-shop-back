package mrlib

import (
    "encoding/json"
    "io"
    "os"
    "print-shop-back/pkg/mrapp"
)

type JsonReader struct {
    logger mrapp.Logger
}

func NewJsonReader(logger mrapp.Logger) *JsonReader {
    return &JsonReader{
        logger: logger,
    }
}

func (jr *JsonReader) Read(fileName string, data any) error {
    jsonFile, err := os.Open(fileName)

    if err != nil {
        return err
    }

    defer jsonFile.Close()

    jr.logger.Info("Successfully Opened %s", fileName)

    contentBytes, err := io.ReadAll(jsonFile)

    if err != nil {
        return err
    }

    jr.logger.Info("File %s successfully readed", fileName)

    err = json.Unmarshal(contentBytes, data)

    return err
}
