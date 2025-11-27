-- Add stream_type and source_id to streams table

ALTER TABLE streams
ADD COLUMN stream_type VARCHAR(50) DEFAULT 'hls',
ADD COLUMN source_id VARCHAR(255);

-- Create index for stream_type
CREATE INDEX idx_streams_stream_type ON streams(stream_type);

COMMENT ON COLUMN streams.stream_type IS 'Type of stream: hls (default), youtube, etc.';
COMMENT ON COLUMN streams.source_id IS 'External ID for the stream source (e.g. YouTube Video ID)';

