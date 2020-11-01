package mapmodels

type Bounds struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}
