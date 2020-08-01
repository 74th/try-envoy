resource "google_container_cluster" "cluster" {
  name                     = "test-gke-cluster"
  location                 = "asia-northeast1"
  initial_node_count       = 1
  remove_default_node_pool = true
}

resource "google_container_node_pool" "node_pool" {
  name               = "node-v1"
  location           = "asia-northeast1"
  node_locations     = ["asia-northeast1-a"]
  cluster            = google_container_cluster.cluster.name
  initial_node_count = 1

  autoscaling {
    min_node_count = 1
    max_node_count = 5
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    preemptible  = true
    machine_type = "n1-standard-2"
    image_type   = "COS"
    disk_type    = "pd-standard"
    disk_size_gb = "30"

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
