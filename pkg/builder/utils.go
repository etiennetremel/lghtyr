package builder

import (
	"bufio"
	"io"
	"sync"

	"github.com/rs/zerolog"
)

// forward logs to zerolog
func logOutput(r io.Reader, log *zerolog.Logger, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		log.Info().Msg(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Error reading std output")
	}
}
