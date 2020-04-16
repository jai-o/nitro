package hack

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"github.com/craftcms/nitro/config"
)

// FindMounts will take a name of a machine and the output of an exec.Command as a slice of bytes
// and return a slice of config mounts that has a source and destination or an error. This is
// used to match if the machine has any mounts. The args passed to multipass are expected to
// return a csv format (e.g. "multipass info machinename --format=csv").
func FindMounts(name string, b []byte) ([]config.Mount, error) {
	var mounts []config.Mount

	rdr := csv.NewReader(bytes.NewReader(b))

	for {
		line, err := rdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if line[0] != name {
			break
		}

		for _, m := range strings.Split(line[12], ";") {
			// since we split on the ;, the last element could be empty
			if len(m) == 0 {
				break
			}

			mount := strings.Split(m, " ")

			mounts = append(mounts, config.Mount{Source: mount[0], Dest: mount[2]})
		}
	}

	return mounts, nil
}
