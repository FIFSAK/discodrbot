package bot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
)

func SaveAudioToFile(audioBuffer [][]int16, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, pcm := range audioBuffer {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, pcm)
		if err != nil {
			return err
		}
		_, err = file.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}

func ConvertPCMToMP3(pcmFilename, mp3Filename string) error {
	cmd := exec.Command("lame", "-r", pcmFilename, mp3Filename)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run() failed with %s\n%s", err, stderr.String())
	}
	return nil
}

func DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
