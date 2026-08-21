-- Create ENUM types
CREATE TYPE user_role AS ENUM ('user', 'seller', 'admin');
CREATE TYPE transaction_type AS ENUM (
  'mining_reward',
  'purchase_discount',
  'subscription_payment',
  'transfer_sent',
  'transfer_received',
  'ads_payment',
  'consensus_validation_reward'
);
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed');
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipped', 'delivered', 'cancelled');
CREATE TYPE ai_request_type AS ENUM ('product_search', 'chat_assistance', 'content_moderation', 'video_recommendation');
CREATE TYPE validation_result AS ENUM ('true', 'false', 'inconclusive');
CREATE TYPE mining_status AS ENUM ('active', 'paused', 'completed');
CREATE TYPE campaign_status AS ENUM ('active', 'paused', 'completed', 'cancelled');

-- Users Table
CREATE TABLE users (
  user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  phone VARCHAR(20) UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  avatar_url TEXT,
  role user_role DEFAULT 'user',
  mining_enabled BOOLEAN DEFAULT FALSE,
  mining_cpu_limit INT DEFAULT 10 CHECK (mining_cpu_limit >= 1 AND mining_cpu_limit <= 10),
  security_settings JSONB DEFAULT '{}'::jsonb,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Wallets Table
CREATE TABLE wallets (
  wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
  balance_tk DECIMAL(18,8) DEFAULT 0 CHECK (balance_tk >= 0),
  frozen_balance_tk DECIMAL(18,8) DEFAULT 0 CHECK (frozen_balance_tk >= 0),
  last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_wallets_user_id ON wallets(user_id);

-- TK Transactions Table
CREATE TABLE tk_transactions (
  transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  transaction_type transaction_type NOT NULL,
  amount_tk DECIMAL(18,8) NOT NULL CHECK (amount_tk > 0),
  description TEXT,
  related_entity_id UUID,
  status transaction_status DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_tk_transactions_user_id ON tk_transactions(user_id);
CREATE INDEX idx_tk_transactions_created_at ON tk_transactions(created_at);
CREATE INDEX idx_tk_transactions_status ON tk_transactions(status);

-- Products Table
CREATE TABLE products (
  product_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  seller_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price_usd DECIMAL(10,2) NOT NULL CHECK (price_usd > 0),
  category VARCHAR(100),
  stock INT DEFAULT 0 CHECK (stock >= 0),
  image_url TEXT,
  ai_search_tags TEXT[] DEFAULT '{}',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE INDEX idx_products_category ON products(category);

-- Orders Table
CREATE TABLE orders (
  order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  buyer_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  seller_id UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  total_price_usd DECIMAL(10,2) NOT NULL CHECK (total_price_usd > 0),
  tk_discount_used DECIMAL(18,8) DEFAULT 0 CHECK (tk_discount_used >= 0),
  final_price_usd DECIMAL(10,2) NOT NULL CHECK (final_price_usd >= 0),
  status order_status DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_orders_buyer_id ON orders(buyer_id);
CREATE INDEX idx_orders_seller_id ON orders(seller_id);
CREATE INDEX idx_orders_status ON orders(status);

-- Order Items Table
CREATE TABLE order_items (
  item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(product_id) ON DELETE RESTRICT,
  quantity INT NOT NULL CHECK (quantity > 0),
  price_per_unit DECIMAL(10,2) NOT NULL CHECK (price_per_unit > 0)
);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);

-- Chats Table
CREATE TABLE chats (
  chat_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  participant_1_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  participant_2_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CHECK (participant_1_id < participant_2_id)
);
CREATE INDEX idx_chats_participants ON chats(participant_1_id, participant_2_id);

-- Messages Table
CREATE TABLE messages (
  message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
  sender_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  encrypted_text BYTEA NOT NULL,
  tk_transfer_amount DECIMAL(18,8),
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);

-- Channels Table
CREATE TABLE channels (
  channel_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_channels_owner_id ON channels(owner_id);

-- Videos Table
CREATE TABLE videos (
  video_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  channel_id UUID NOT NULL REFERENCES channels(channel_id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  file_path VARCHAR(512) NOT NULL,
  file_hash VARCHAR(64),
  duration_seconds INT,
  view_count INT DEFAULT 0,
  p2p_enabled BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_videos_channel_id ON videos(channel_id);

-- Video Views Table
CREATE TABLE video_views (
  view_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  video_id UUID NOT NULL REFERENCES videos(video_id) ON DELETE CASCADE,
  viewer_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  watched_seconds INT DEFAULT 0,
  p2p_contributed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_video_views_video_id ON video_views(video_id);
CREATE INDEX idx_video_views_viewer_id ON video_views(viewer_id);

-- AI Requests Table
CREATE TABLE ai_requests (
  request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  request_type ai_request_type NOT NULL,
  query_text TEXT NOT NULL,
  response_text TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_ai_requests_user_id ON ai_requests(user_id);

-- Consensus Validations Table
CREATE TABLE consensus_validations (
  validation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ai_request_id UUID NOT NULL REFERENCES ai_requests(request_id) ON DELETE CASCADE,
  validator_node_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  fact_claim TEXT NOT NULL,
  validation_result validation_result NOT NULL,
  confidence_score DECIMAL(3,2) CHECK (confidence_score >= 0 AND confidence_score <= 1),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_consensus_validations_ai_request_id ON consensus_validations(ai_request_id);

-- Mining Sessions Table
CREATE TABLE mining_sessions (
  session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  started_at TIMESTAMP NOT NULL,
  paused_at TIMESTAMP,
  duration_minutes INT DEFAULT 0,
  facts_verified INT DEFAULT 0,
  reward_tk_earned DECIMAL(18,8) DEFAULT 0 CHECK (reward_tk_earned >= 0),
  status mining_status DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_mining_sessions_user_id ON mining_sessions(user_id);
CREATE INDEX idx_mining_sessions_status ON mining_sessions(status);

-- Ad Campaigns Table
CREATE TABLE ad_campaigns (
  campaign_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  creator_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  banner_url TEXT,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  budget_tk DECIMAL(18,8) NOT NULL CHECK (budget_tk > 0),
  spent_tk DECIMAL(18,8) DEFAULT 0 CHECK (spent_tk >= 0),
  status campaign_status DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CHECK (start_date <= end_date)
);
CREATE INDEX idx_ad_campaigns_creator_id ON ad_campaigns(creator_id);

-- Ad Impressions Table
CREATE TABLE ad_impressions (
  impression_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  campaign_id UUID NOT NULL REFERENCES ad_campaigns(campaign_id) ON DELETE CASCADE,
  viewer_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_ad_impressions_campaign_id ON ad_impressions(campaign_id);

-- Seed test data
INSERT INTO users (email, phone, password_hash, role, mining_enabled) VALUES
  ('user1@techno.re', '+1234567890', '$2a$12$example_hash_user1', 'user', true),
  ('seller1@techno.re', '+1234567891', '$2a$12$example_hash_seller1', 'seller', false),
  ('admin@techno.re', '+1234567892', '$2a$12$example_hash_admin', 'admin', false);

INSERT INTO wallets (user_id, balance_tk, frozen_balance_tk)
SELECT user_id, 1000.00000000, 0 FROM users WHERE email = 'user1@techno.re';

INSERT INTO wallets (user_id, balance_tk, frozen_balance_tk)
SELECT user_id, 500.00000000, 0 FROM users WHERE email = 'seller1@techno.re';

INSERT INTO wallets (user_id, balance_tk, frozen_balance_tk)
SELECT user_id, 0, 0 FROM users WHERE email = 'admin@techno.re';
