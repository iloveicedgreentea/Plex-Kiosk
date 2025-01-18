import { useEffect, useState } from "react"
import { Libraries } from "./types"
import { LibrarySection } from "./components/LibrarySection"

function App() {
  const [libraries, setLibraries] = useState<Libraries>({});
  const [lastUpdated, setLastUpdated] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('/api/');
        if (!response.ok) {
          throw new Error('Failed to fetch library data');
        }
        const data = await response.json();
        setLibraries(data.libraries);
        setLastUpdated(data.lastUpdated);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#1a1a1a] text-white flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-[#1a1a1a] text-white flex items-center justify-center">
        <div className="text-red-500 text-xl">{error}</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#1a1a1a] text-white p-6">
      <div className="max-w-[1400px] mx-auto">
        <h1 className="text-3xl font-bold text-[#e5a00d] mb-2">Plex Library</h1>
        <div className="text-[#b3b3b3] text-sm mb-8">
          Last updated: {lastUpdated}
        </div>

        {Object.entries(libraries).map(([name, items]) => (
          <LibrarySection key={name} name={name} items={items} />
        ))}
      </div>
    </div>
  );
}

export default App;