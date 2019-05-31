package dashboard

import "github.com/LinMAD/BitAccretion/model"

// getMaxRequestValue max value from vertices
func getMaxRequestValue(g *model.Graph) (max int) {
	allVertices := g.GetAllVertices()

	for i := 0; i < len(allVertices); i++ {
		if max < allVertices[i].Metric.RequestCount {
			max = allVertices[i].Metric.RequestCount
		}
	}

	return
}
