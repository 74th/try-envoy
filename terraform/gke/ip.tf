resource "google_compute_global_address" "lb_ip" {
  name = "try-envoy-ip"

}

output "lb_ip" {
  value = "${google_compute_global_address.lb_ip.address}"
}
