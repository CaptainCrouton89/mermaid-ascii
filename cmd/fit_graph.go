package cmd

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

type graphFitPlan struct {
	paddingX          int
	paddingY          int
	boxBorderPadding  int
	graphDirection    string
	labelWrapWidth    int
	edgeLabelPolicy   string
	edgeLabelMaxWidth int
	maxNodesPerLevel  int
}

func fitGraphToWidth(properties *graphProperties, config *diagram.Config) string {
	basePlan := graphFitPlan{
		paddingX:          properties.paddingX,
		paddingY:          Min(properties.paddingY, 3),
		boxBorderPadding:  properties.boxBorderPadding,
		graphDirection:    properties.graphDirection,
		labelWrapWidth:    properties.labelWrapWidth,
		edgeLabelPolicy:   properties.edgeLabelPolicy,
		edgeLabelMaxWidth: properties.edgeLabelMaxWidth,
	}

	plans := graphFitPlans(basePlan, config.MaxWidth)
	bestOutput := ""
	bestWidth := 0
	// Track best non-lossy result (preserves edge labels)
	bestLosslessOutput := ""
	bestLosslessWidth := 0
	for idx, plan := range plans {
		candidate := applyGraphFitPlan(properties, plan)
		output := drawMap(candidate)
		width := maxOutputLineWidth(output)
		lossy := plan.edgeLabelPolicy == diagram.EdgeLabelPolicyDrop ||
			plan.edgeLabelPolicy == diagram.EdgeLabelPolicyEllipsis
		if idx == 0 || width < bestWidth {
			bestWidth = width
			bestOutput = output
		}
		if !lossy && (bestLosslessOutput == "" || width < bestLosslessWidth) {
			bestLosslessWidth = width
			bestLosslessOutput = output
		}
		if width <= config.MaxWidth {
			return output
		}
	}

	// Prefer lossless result unless lossy is meaningfully narrower
	if bestLosslessOutput != "" && (bestLosslessWidth <= bestWidth+10) {
		return bestLosslessOutput
	}
	return bestOutput
}

func applyGraphFitPlan(base *graphProperties, plan graphFitPlan) *graphProperties {
	candidate := *base
	candidate.paddingX = plan.paddingX
	candidate.paddingY = plan.paddingY
	candidate.boxBorderPadding = plan.boxBorderPadding
	candidate.graphDirection = plan.graphDirection
	candidate.labelWrapWidth = plan.labelWrapWidth
	candidate.edgeLabelPolicy = plan.edgeLabelPolicy
	candidate.edgeLabelMaxWidth = plan.edgeLabelMaxWidth
	candidate.maxNodesPerLevel = plan.maxNodesPerLevel
	return &candidate
}

func graphFitPlans(base graphFitPlan, maxWidth int) []graphFitPlan {
	plans := []graphFitPlan{}
	seen := map[string]struct{}{}
	addPlan := func(plan graphFitPlan) {
		key := graphFitPlanKey(plan)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		plans = append(plans, plan)
	}

	addPlan(base)

	compact := base
	compact.paddingX = Min(base.paddingX, 2)
	compact.paddingY = Min(base.paddingY, 3)
	compact.boxBorderPadding = Min(base.boxBorderPadding, 1)
	addPlan(compact)

	tight := compact
	tight.paddingX = Min(tight.paddingX, 1)
	tight.paddingY = Min(tight.paddingY, 3)
	tight.boxBorderPadding = 0
	addPlan(tight)

	wrap := base
	wrap.labelWrapWidth = reduceWrapWidth(base.labelWrapWidth, labelWrapWidthFor(maxWidth, base.boxBorderPadding))
	addPlan(wrap)

	wrapCompact := compact
	wrapCompact.labelWrapWidth = reduceWrapWidth(compact.labelWrapWidth, labelWrapWidthFor(maxWidth, compact.boxBorderPadding))
	addPlan(wrapCompact)

	wrapTight := tight
	wrapTight.labelWrapWidth = reduceWrapWidth(tight.labelWrapWidth, labelWrapWidthFor(maxWidth, tight.boxBorderPadding))
	addPlan(wrapTight)

	flippedDirection := flipGraphDirection(base.graphDirection)
	if flippedDirection != "" {
		flipped := base
		flipped.graphDirection = flippedDirection
		addPlan(flipped)

		flippedWrap := flipped
		flippedWrap.labelWrapWidth = reduceWrapWidth(flipped.labelWrapWidth, labelWrapWidthFor(maxWidth, flipped.boxBorderPadding))
		addPlan(flippedWrap)

		flippedCompact := compact
		flippedCompact.graphDirection = flippedDirection
		flippedCompact.labelWrapWidth = reduceWrapWidth(flippedCompact.labelWrapWidth, labelWrapWidthFor(maxWidth, flippedCompact.boxBorderPadding))
		addPlan(flippedCompact)
	}

	// Reflow: limit nodes per level to force wrapping.
	// Reflow needs enough horizontal padding for edge routing columns
	// between side-by-side nodes — without it edges cut through labels.
	for _, limit := range []int{4, 3, 2} {
		reflow := compact
		reflow.maxNodesPerLevel = limit
		reflow.paddingX = Max(compact.paddingX, 5)
		addPlan(reflow)

		reflowTight := compact
		reflowTight.maxNodesPerLevel = limit
		reflowTight.paddingX = Max(compact.paddingX, 3)
		reflowTight.boxBorderPadding = 0
		addPlan(reflowTight)

		reflowWrap := wrapCompact
		reflowWrap.maxNodesPerLevel = limit
		reflowWrap.paddingX = Max(wrapCompact.paddingX, 5)
		addPlan(reflowWrap)
	}

	ellipsis := wrapCompact
	ellipsis.edgeLabelPolicy = diagram.EdgeLabelPolicyEllipsis
	ellipsis.edgeLabelMaxWidth = edgeLabelMaxWidthFor(maxWidth)
	addPlan(ellipsis)

	drop := ellipsis
	drop.edgeLabelPolicy = diagram.EdgeLabelPolicyDrop
	addPlan(drop)

	return plans
}

func graphFitPlanKey(plan graphFitPlan) string {
	return fmt.Sprintf("%d:%d:%d:%s:%d:%s:%d:%d",
		plan.paddingX,
		plan.paddingY,
		plan.boxBorderPadding,
		plan.graphDirection,
		plan.labelWrapWidth,
		plan.edgeLabelPolicy,
		plan.edgeLabelMaxWidth,
		plan.maxNodesPerLevel,
	)
}

func labelWrapWidthFor(maxWidth, boxBorderPadding int) int {
	if maxWidth <= 0 {
		return 0
	}
	available := maxWidth - (2*boxBorderPadding + 2)
	if available < 1 {
		return 1
	}
	return available
}

func edgeLabelMaxWidthFor(maxWidth int) int {
	if maxWidth <= 0 {
		return 0
	}
	if maxWidth <= 3 {
		return maxWidth
	}
	return maxWidth - 4
}

func reduceWrapWidth(current, target int) int {
	if current <= 0 || target < current {
		return target
	}
	return current
}

func flipGraphDirection(direction string) string {
	switch direction {
	case "LR":
		return "TD"
	case "TD":
		return "LR"
	default:
		return ""
	}
}

func maxOutputLineWidth(output string) int {
	lines := strings.Split(output, "\n")
	max := 0
	for _, line := range lines {
		if len(line) > max {
			max = len(line)
		}
	}
	return max
}
