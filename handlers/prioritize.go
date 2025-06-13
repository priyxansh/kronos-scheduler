package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	corev1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kube-scheduler/extender/v1"
)

// Prioritize handles the /prioritize extender endpoint
func Prioritize(c *fiber.Ctx) error {
	var args schedulerapi.ExtenderArgs
	if err := c.BodyParser(&args); err != nil {
		log.Printf("Failed to parse ExtenderArgs: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to decode request"})
	}

	jobType := args.Pod.Labels["jobType"]
	if jobType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing jobType label"})
	}

	var priorityList []schedulerapi.HostPriority
	for _, node := range args.Nodes.Items {
		score := calculateScore(node, jobType)
		log.Printf("Scored node %s for jobType=%s → %d", node.Name, jobType, score)

		priorityList = append(priorityList, schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		})
	}

	return c.JSON(priorityList)
}

// calculateScore returns a node score (0–100) based on job type and simplified queue model
func calculateScore(node corev1.Node, jobType string) int64 {
	// Use Allocatable CPU as a proxy for system load (for demo purposes)
	cpuAlloc := node.Status.Allocatable.Cpu().MilliValue()
	rho := float64(cpuAlloc) / 1000.0 // Normalize to range like 0–2

	// Placeholder queue model values (can be dynamically tuned later)
	lambda := 0.1 // arrival rate
	varS := 1.0   // service time variance

	if jobType == "long" {
		// M/G/1: Wq = (λ × Var(S)) / (2 × (1 – ρ))
		if rho >= 1.0 {
			return 10 // prevent negative or undefined wait time when overloaded
		}
		wq := (lambda * varS) / (2 * (1 - rho))
		score := int64(100 - wq*10)
		if score < 0 {
			return 0
		}
		return score
	}

	// M/M/1 for short jobs → lighter penalty
	score := int64(100 - rho*10)
	if score < 0 {
		return 0
	}
	return score
}
