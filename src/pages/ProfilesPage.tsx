import { useEffect, useState } from "react";
import { api } from "@/services/api";
import type { Profile } from "@/types";
import "./ProfilesPage.css";

interface ProfilesPageProps {
  onSelectProfile: () => void;
}

export default function ProfilesPage({ onSelectProfile }: ProfilesPageProps) {
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddProfile, setShowAddProfile] = useState(false);
  const [newProfileName, setNewProfileName] = useState("");
  const [isKidsProfile, setIsKidsProfile] = useState(false);

  useEffect(() => {
    loadProfiles();
  }, []);

  const loadProfiles = async () => {
    try {
      const data = await api.getProfiles();
      setProfiles(data);
    } catch (error) {
      console.error("Failed to load profiles:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newProfileName.trim()) return;

    try {
      const newProfile = await api.createProfile(newProfileName, isKidsProfile);
      setProfiles([...profiles, newProfile]);
      setNewProfileName("");
      setIsKidsProfile(false);
      setShowAddProfile(false);
    } catch (error) {
      console.error("Failed to create profile:", error);
    }
  };

  if (loading) {
    return <div className="loading">Loading profiles...</div>;
  }

  return (
    <div className="profiles-page">
      <div className="profiles-container">
        <h1>Who's watching?</h1>

        <div className="profiles-grid">
          {profiles.map((profile) => (
            <div
              key={profile.id}
              className="profile-card"
              onClick={onSelectProfile}
            >
              <div className="profile-avatar">
                <img src={profile.avatar_url} alt={profile.name} />
              </div>
              <h3>{profile.name}</h3>
            </div>
          ))}

          {!showAddProfile && (
            <div
              className="profile-card add-profile"
              onClick={() => setShowAddProfile(true)}
            >
              <div className="profile-avatar">
                <span className="add-icon">+</span>
              </div>
              <h3>Add Profile</h3>
            </div>
          )}
        </div>

        {showAddProfile && (
          <div className="add-profile-form">
            <h2>Add Profile</h2>
            <form onSubmit={handleCreateProfile}>
              <input
                type="text"
                placeholder="Profile Name"
                value={newProfileName}
                onChange={(e) => setNewProfileName(e.target.value)}
                required
                autoFocus
              />

              <label className="checkbox-label">
                <input
                  type="checkbox"
                  checked={isKidsProfile}
                  onChange={(e) => setIsKidsProfile(e.target.checked)}
                />
                <span>Kids profile (shows appropriate content only)</span>
              </label>

              <div className="form-buttons">
                <button type="submit" className="submit-button">
                  Create
                </button>
                <button
                  type="button"
                  className="cancel-button"
                  onClick={() => {
                    setShowAddProfile(false);
                    setNewProfileName("");
                    setIsKidsProfile(false);
                  }}
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}
      </div>
    </div>
  );
}
