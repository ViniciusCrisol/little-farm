export async function createFarm(mapWidth, mapHeight) {
	const response = await fetch("/api/farms", {
		body: JSON.stringify({
			map_width: mapWidth,
			map_height: mapHeight,
		}),
		method: "POST",
	});
	return parseResponse(await response.json());
}

function parseResponse(response) {
	return {
		farmId: response.farm_id,
		mapWidth: response.map_width,
		mapHeight: response.map_height,
		activeFor: response.active_for,
		resources: {
			corn: response.resources.corn,
			seeds: response.resources.seeds,
			moneyInCents: response.resources.money_in_cents,
		},
		activeElements: response.active_elements.map((e) => ({
			id: e.id,
			kind: e.kind,
			xPosition: e.x_position,
			yPosition: e.y_position,
		})),
		passiveElements: response.passive_elements.map((e) => ({
			id: e.id,
			kind: e.kind,
		})),
		createdAt: response.created_at,
		updatedAt: response.updated_at,
	};
}
