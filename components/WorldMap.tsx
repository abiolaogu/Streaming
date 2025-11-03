import React from 'react';

type Region = 'all' | 'americas' | 'emea' | 'apac';

interface WorldMapProps {
  activeRegion: Region;
}

const WorldMap: React.FC<WorldMapProps> = ({ activeRegion }) => {
  const getRegionClass = (region: Region) => {
    return `transition-colors duration-300 ${
      activeRegion === region ? 'fill-brand-accent/40 dark:fill-brand-accent/30' : 'fill-brand-border'
    }`;
  };

  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 1000 500"
      className="w-full h-full"
      aria-label="World map with selectable regions"
    >
      {/* Americas */}
      <path
        className={getRegionClass('americas')}
        d="M205 65 l-25 15 v15 l10 10 h15 l10 10 v15 l-10 10 h-15 l-10 10 v15 l10 10 h15 l10 10 v15 l-10 10 h-15 l-10 10 v15 l10 10 h15 l10 10 v15 l-10 10 h-15 l-10 10 v15 l10 10 h15 l10 10 v15 h30 l15 -10 v-15 l-10 -10 h-15 l-10 -10 v-15 l10 -10 h15 l10 -10 v-15 l-10 -10 h-15 l-10 -10 v-15 l10 -10 h15 l10 -10 v-15 h-30 l-15 10 v15 z M 300 250 l-20 30 l-10 20 v20 l10 10 h20 l10 -10 v-20 l-10 -20 h-20 l-10 10 v20 h20 l10 -10 v-20 z M150,300 l0,15 -10,15 -10,0 0,-15 10,-15 10,0z m-30-200 -15,0 -15,15 -30,90 0,30 15,30 15,0 15,-15 0,-15 -15,0 0,-15 15,-15 15,0 0,-30 -15,-30 -15,0z"
      />
      
      {/* EMEA (Europe, Middle East, Africa) */}
      <path
        className={getRegionClass('emea')}
        d="M450 120 l-10 20 h-10 l-10 20 v20 l10 10 h10 l10 -10 v-20 h10 l10 -10 v-20 h-10 z M 480 150 v80 l-20 20 h-20 l-10 30 v40 l10 20 h20 l20 10 v30 l-10 10 h-30 l-20 20 v30 l30 20 h40 l20-20 v-180 l-10-10 h-20z M 550 100 l15,0 15,15 0,30 -15,15 -15,0z m30-30 0,15 15,15 15,0 0,-15 -15,-15 -15,0z m-60 0 -15,0 -15,15 0,30 15,15 15,0z"
      />

      {/* APAC (Asia-Pacific) */}
      <path
        className={getRegionClass('apac')}
        d="M650 100 l20,0 20,20 0,30 -20,20 -20,0z m100 50 20,0 20,20 0,30 -20,20 -20,0z M 850 350 l0,20 20,20 30,0 0,-20 -20,-20 -30,0z M700,400 l-15,15 -15,0 0,-15 15,-15 15,0z M900,450 l0,15 -15,15 -15,0 0,-15 15,-15 15,0z M 950 200 l-30 30 v30 l15 15 h30 l15-15 v-30 l-15-15 h-30z M620 150 l100 50 v50 l-50 50 h-50 l-20-20 v-80 l20-30 z M750,280 l30,0 20,20 0,50 -20,20 -30,0z"
      />
    </svg>
  );
};

export default WorldMap;