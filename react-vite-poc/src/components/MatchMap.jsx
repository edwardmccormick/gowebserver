import { useEffect, useRef } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

function MatchMap({ people, User }) {
  const mapContainer = useRef(null);
  const map = useRef(null);

  useEffect(() => {
    if (!mapContainer.current || !User) return;

    // Set Mapbox access token from environment variable
    mapboxgl.accessToken = import.meta.env.VITE_MAPBOX_API_KEY;

    // Initialize map
    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: 'mapbox://styles/mapbox/streets-v12',
      center: [User.long, User.lat],
      zoom: 10
    }).addControl(new mapboxgl.NavigationControl(), 'top-right')
      .addControl(new mapboxgl.FullscreenControl(), 'top-right')
      .addControl(new mapboxgl.GeolocateControl({
        positionOptions: {
          enableHighAccuracy: true
        },
        trackUserLocation: true,
        showUserHeading: true
      }), 'top-right');

    // Add user marker
    new mapboxgl.Marker({ color: 'red' })
      .setLngLat([User.long, User.lat])
      .setPopup(new mapboxgl.Popup().setHTML(`<h6>You</h6><p>${User.name}</p>`))
      .addTo(map.current);

    // Add match markers
    people.forEach((person) => {
      new mapboxgl.Marker({ color: 'blue' })
        .setLngLat([person.long, person.lat])
        .setPopup(new mapboxgl.Popup().setHTML(`
          <div style="text-align: center;">
            <img src="${person.profile.url || '/profile.svg'}" style="width: 50px; height: 50px; border-radius: 50%; margin-bottom: 8px;" />
            <h6>${person.name}</h6>
            <p>${person.motto}</p>
            <small>${person.distance} miles away</small>
          </div>
        `))
        .addTo(map.current);
    });

    return () => map.current?.remove();
  }, [people, User]);

  return <div ref={mapContainer} style={{ width: '100%', height: '500px' }} />;
}

export default MatchMap;