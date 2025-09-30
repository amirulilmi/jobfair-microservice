-- Create company_media table for better media management
CREATE TABLE IF NOT EXISTS company_media (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    
    -- Media Information
    media_type VARCHAR(50) NOT NULL CHECK (media_type IN ('logo', 'banner', 'video', 'gallery', 'document')),
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_size BIGINT, -- in bytes
    mime_type VARCHAR(100),
    
    -- Media Metadata
    title VARCHAR(255),
    description TEXT,
    alt_text VARCHAR(255),
    display_order INTEGER DEFAULT 0,
    
    -- Processing Status
    is_processed BOOLEAN DEFAULT TRUE,
    is_public BOOLEAN DEFAULT TRUE,
    thumbnail_url VARCHAR(500),
    
    -- Analytics
    view_count INTEGER DEFAULT 0,
    download_count INTEGER DEFAULT 0,
    
    -- Timestamps
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Foreign Key
    CONSTRAINT fk_company_media_company_id 
        FOREIGN KEY (company_id) 
        REFERENCES companies(id) 
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_company_media_company_id ON company_media(company_id);
CREATE INDEX idx_company_media_type ON company_media(media_type);
CREATE INDEX idx_company_media_display_order ON company_media(display_order);
CREATE INDEX idx_company_media_is_public ON company_media(is_public);
CREATE INDEX idx_company_media_deleted_at ON company_media(deleted_at);

-- Create trigger for updated_at
CREATE TRIGGER update_company_media_updated_at BEFORE UPDATE ON company_media
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE company_media IS 'Media files for company virtual booth';
COMMENT ON COLUMN company_media.media_type IS 'Type: logo, banner, video, gallery, or document';
COMMENT ON COLUMN company_media.display_order IS 'Order for displaying media (0 is first)';
