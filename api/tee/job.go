package tee

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/masa-finance/tee-worker/api/types"
	"github.com/masa-finance/tee-worker/pkg/tee"
	"golang.org/x/exp/rand"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		// TODO: Move xcrypt from indexer to tee-types, and use RandomString here (although we'll need a different alpahbet)
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GenerateJobSignature generates a signature for the job.
func GenerateJobSignature(job *types.Job) (string, error) {
	dat, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	checksum := sha256.New()
	checksum.Write(dat)

	job.Nonce = fmt.Sprintf("%s-%s", string(checksum.Sum(nil)), randStringRunes(99))

	dat, err = json.Marshal(job)
	if err != nil {
		return "", err
	}

	return tee.Seal(dat)
}

// SealJobResult seals a job result with the job's nonce.
func SealJobResult(jr *types.JobResult) (string, error) {
	return tee.SealWithKey(jr.Job.Nonce, jr.Data)
}

// DecryptJob decrypts the job request.
func DecryptJob(jobRequest *types.JobRequest) (*types.Job, error) {
	dat, err := tee.Unseal(jobRequest.EncryptedJob)
	if err != nil {
		return nil, err
	}

	job := types.Job{}
	if err := json.Unmarshal(dat, &job); err != nil {
		return nil, err
	}

	return &job, nil
}
