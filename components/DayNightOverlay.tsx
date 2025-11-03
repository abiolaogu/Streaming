
import React, { useState, useEffect } from 'react';

const DayNightOverlay: React.FC = () => {
    const [gradient, setGradient] = useState('');

    useEffect(() => {
        const updateOverlay = () => {
            const now = new Date();
            // Get day of the year (1-366)
            const startOfYear = new Date(now.getFullYear(), 0, 0);
            const diff = (now.getTime() - startOfYear.getTime()) + ((startOfYear.getTimezoneOffset() - now.getTimezoneOffset()) * 60 * 1000);
            const oneDay = 1000 * 60 * 60 * 24;
            const dayOfYear = Math.floor(diff / oneDay);
            
            const utcHours = now.getUTCHours() + now.getUTCMinutes() / 60;

            // Calculate sun's declination for seasonal tilt effect
            const declination = -23.45 * Math.cos((2 * Math.PI / 365) * (dayOfYear + 10));
            const gradientAngle = 90 - declination;

            const dayFraction = utcHours / 24;
            const centerPercent = dayFraction * 100;

            const nightColor = 'rgba(13, 17, 23, 0.4)'; // A semi-transparent dark color matching brand-bg
            const dayColor = 'rgba(255, 255, 255, 0.0)'; // Transparent for day

            const bandWidthPercent = 50; // The width of the "day" band
            const transitionWidthPercent = 7; // The width of the sunrise/sunset transition

            let startDay = centerPercent - (bandWidthPercent / 2);
            let endDay = centerPercent + (bandWidthPercent / 2);

            let cssGradient;

            if (startDay < 0) { // Wraps around the left edge
                const wrappedStart = startDay + 100;
                cssGradient = `linear-gradient(${gradientAngle}deg, 
                    ${dayColor} ${endDay}%, 
                    ${nightColor} ${endDay + transitionWidthPercent}%, 
                    ${nightColor} ${wrappedStart - transitionWidthPercent}%, 
                    ${dayColor} ${wrappedStart}%
                )`;
            } else if (endDay > 100) { // Wraps around the right edge
                const wrappedEnd = endDay - 100;
                cssGradient = `linear-gradient(${gradientAngle}deg, 
                    ${nightColor} ${wrappedEnd - transitionWidthPercent}%, 
                    ${dayColor} ${wrappedEnd}%, 
                    ${dayColor} ${startDay}%, 
                    ${nightColor} ${startDay + transitionWidthPercent}%
                )`;
            } else { // Day is fully visible in one segment
                cssGradient = `linear-gradient(${gradientAngle}deg, 
                    ${nightColor}, 
                    ${nightColor} ${startDay - transitionWidthPercent}%, 
                    ${dayColor} ${startDay}%, 
                    ${dayColor} ${endDay}%, 
                    ${nightColor} ${endDay + transitionWidthPercent}%, 
                    ${nightColor}
                )`;
            }
            
            setGradient(cssGradient);
        };

        updateOverlay();
        const intervalId = setInterval(updateOverlay, 60000); // Update every minute

        return () => clearInterval(intervalId);
    }, []);

    return <div className="absolute inset-0 w-full h-full pointer-events-none" style={{ background: gradient }} />;
};

export default DayNightOverlay;
