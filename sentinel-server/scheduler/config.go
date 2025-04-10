package scheduler

type Config struct {
    SubfinderInterval int // hours
    HttpxInterval     int // hours
    DnsxInterval      int // hours
}

func NewDefaultConfig() Config {
    return Config{
        SubfinderInterval: 24, // run every 24 hours
        HttpxInterval:     12, // run every 12 hours
        DnsxInterval:      6,  // run every 6 hours
    }
}