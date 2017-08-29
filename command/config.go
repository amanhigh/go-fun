package commander

import "os"

/* PSSH */
const CLUSTER_PATH = "/tmp/clusters"
const OUTPUT_PATH = "/tmp/output"
const ERROR_PATH = "/tmp/error"
const CONSOLE_FILE = CLUSTER_PATH + "/console.txt"

const RELEASE_FILE = "~/Documents/release.txt"
const DEFAULT_PERM = os.FileMode(0755)

/* Date/Time */
const LAYOUT = "Jan 2, 2006 at 3:04pm (MST)"

