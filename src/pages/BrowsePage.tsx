import { useCallback, useEffect, useState } from "react";
import NavBar from "@/components/NavBar";
import { api } from "@/services/api";
import type { Content } from "@/types";
import type { Page } from "@/types/navigation";
import "./BrowsePage.css";

interface BrowsePageProps {
  onWatchContent: (contentId: string) => void;
  onNavigate: (page: Page) => void;
  onLogout: () => void;
}

export default function BrowsePage({
  onWatchContent,
  onNavigate,
  onLogout,
}: BrowsePageProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<Content[]>([]);
  const [trending, setTrending] = useState<Content[]>([]);
  const [loading, setLoading] = useState(false);

  const loadTrending = useCallback(async () => {
    try {
      const data = await api.getTrending();
      setTrending(data);
    } catch (error) {
      console.error("Failed to load trending:", error);
    }
  }, []);

  const handleSearch = useCallback(async () => {
    if (!searchQuery.trim()) return;

    setLoading(true);
    try {
      const results = await api.searchContent(searchQuery);
      setSearchResults(results);
    } catch (error) {
      console.error("Search failed:", error);
    } finally {
      setLoading(false);
    }
  }, [searchQuery]);

  useEffect(() => {
    loadTrending();
  }, [loadTrending]);

  useEffect(() => {
    if (searchQuery.trim()) {
      handleSearch();
    } else {
      setSearchResults([]);
    }
  }, [searchQuery, handleSearch]);

  const displayContent = searchQuery.trim() ? searchResults : trending;

  return (
    <div className="browse-page">
      <NavBar onNavigate={onNavigate} onLogout={onLogout} />

      <div className="browse-content">
        <div className="search-section">
          <input
            type="text"
            placeholder="Search for movies, series, documentaries..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-input"
          />
        </div>

        <h2 className="section-title">
          {searchQuery.trim() ? "Search Results" : "Trending Now"}
        </h2>

        {loading ? (
          <div className="loading">Searching...</div>
        ) : (
          <div className="content-grid">
            {displayContent.map((content) => (
              <div
                key={content.id}
                className="content-card"
                onClick={() => onWatchContent(content.id)}
              >
                <img src={content.thumbnail_url} alt={content.title} />
                <div className="card-overlay">
                  <h3>{content.title}</h3>
                  <div className="card-info">
                    <span>★ {content.rating.toFixed(1)}</span>
                    <span>{content.type}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {!loading && displayContent.length === 0 && (
          <div className="no-results">
            {searchQuery.trim()
              ? "No results found"
              : "No trending content available"}
          </div>
        )}
      </div>
    </div>
  );
}
