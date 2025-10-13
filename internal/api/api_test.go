package api_test

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	teetypes "github.com/masa-finance/tee-worker/api/types"
	. "github.com/masa-finance/tee-worker/internal/api"
	"github.com/masa-finance/tee-worker/internal/config"
	"github.com/masa-finance/tee-worker/pkg/client"
)

var _ = Describe("API", func() {

	var (
		clientInstance *client.Client
		ctx            context.Context
		cancel         context.CancelFunc
	)

	BeforeEach(func() {
		// Start the server
		ctx, cancel = context.WithCancel(context.Background())

		// Start the server
		os.Setenv("LOG_LEVEL", "debug")
		logrus.SetLevel(logrus.DebugLevel)
		go func() {
			logrus.SetLevel(logrus.DebugLevel)
			Start(ctx, "127.0.0.1:40912", "", true, config.JobConfiguration{})
		}()

		// Wait for the server to start
		Eventually(func() error {

			c, err := client.NewClient("http://localhost:40912")
			if err != nil {
				return err
			}

			signature, err := c.CreateJobSignature(teetypes.Job{
				Type:      teetypes.WebJob,
				Arguments: map[string]interface{}{},
			})
			if err != nil {
				return err
			}

			// Check if the job signature is non-empty (indicates server is ready)
			if signature == "" {
				return fmt.Errorf("job signature is empty, server not ready")
			}

			return nil // Success: signature is non-empty
		}, 10*time.Second).Should(Succeed())

		// Initialize the client
		clientInstance, _ = client.NewClient("http://localhost:40912")
	})

	AfterEach(func() {
		// Stop the server
		cancel()
	})

	It("should submit a job and get the correct result", func() {
		// Step 1: Create the job request
		// we use TikTok transcription here as it's supported by all workers without any unique config
		job := teetypes.Job{
			Type: teetypes.TiktokJob,
			Arguments: map[string]interface{}{
				"type":      "transcription",
				"video_url": "https://www.tiktok.com/@theblockrunner.com/video/7227579907361066282",
			},
		}
		// Step 2: Get a Job signature
		jobSignature, err := clientInstance.CreateJobSignature(job)
		Expect(err).NotTo(HaveOccurred())
		Expect(jobSignature).NotTo(BeEmpty())

		// Step 3: Submit the job
		jobResult, err := clientInstance.SubmitJob(jobSignature)
		Expect(err).NotTo(HaveOccurred())
		Expect(jobResult).NotTo(BeNil())

		// Step 4: Wait for the job result
		encryptedResult, err := jobResult.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(encryptedResult).NotTo(BeEmpty())

		// Step 5: Decrypt the result
		decryptedResult, err := clientInstance.Decrypt(jobSignature, encryptedResult)
		Expect(err).NotTo(HaveOccurred())

		Expect(decryptedResult).NotTo(BeEmpty())

		result, err := jobResult.GetDecrypted(jobSignature)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).NotTo(BeEmpty())
	})

	It("bubble up errors", func() {
		// Step 1: Create the job request
		job := teetypes.Job{
			Type: "not-existing scraper",
			Arguments: map[string]interface{}{
				"url": "google",
			},
		}

		// Step 2: Get a Job signature
		jobSignature, err := clientInstance.CreateJobSignature(job)
		Expect(err).NotTo(HaveOccurred())
		Expect(jobSignature).NotTo(BeEmpty())

		// Step 3: Submit the job
		jobResult, err := clientInstance.SubmitJob(jobSignature)
		Expect(err).NotTo(HaveOccurred())
		Expect(jobResult).NotTo(BeNil())
		Expect(jobResult.UUID).NotTo(BeEmpty())

		jobResult.SetMaxRetries(10)

		// Step 4: Wait for the job result (should fail)
		encryptedResult, err := jobResult.Get()
		Expect(err).To(HaveOccurred())
		Expect(encryptedResult).To(BeEmpty())
	})
})
