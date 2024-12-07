package query

type QueryMetadata struct {
    QueryID     string                 `json:"query_id"`
    QueryType   string                 `json:"query_type"`
    StartTime   time.Time              `json:"start_time"`
    EndTime     time.Time              `json:"end_time"`
    Duration    time.Duration          `json:"duration"`
    CacheHit    bool                   `json:"cache_hit"`
    CacheKey    string                 `json:"cache_key,omitempty"`
    TraceID     string                 `json:"trace_id"`
    SpanID      string                 `json:"span_id"`
    UserID      string                 `json:"user_id,omitempty"`
    IPAddress   string                 `json:"ip_address,omitempty"`
    UserAgent   string                 `json:"user_agent,omitempty"`
    Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

func NewQueryMetadata() *QueryMetadata {
    return &QueryMetadata{
        QueryID:    uuid.New().String(),
        StartTime:  time.Now(),
        Attributes: make(map[string]interface{}),
    }
}

func (m *QueryMetadata) Complete() {
    m.EndTime = time.Now()
    m.Duration = m.EndTime.Sub(m.StartTime)
} 