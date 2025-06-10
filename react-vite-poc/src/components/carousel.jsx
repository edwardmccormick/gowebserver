import { useState, useEffect } from 'react';
import Carousel from 'react-bootstrap/Carousel';
import Image from 'react-bootstrap/Image';

function ControlledCarousel() {
//   const [index, setIndex] = useState(0);

//   const handleSelect = (selectedIndex) => {
//     setIndex(selectedIndex);
//   };

    const [photos, setPhotos] = useState([]);
    const [loading, setLoading] = useState(true);
  
    useEffect(() => {
      fetch('http://localhost:8080/photos/1')
        .then((res) => res.json())
        .then((data) => {
          setPhotos(data);
          setLoading(false);
        })
        .catch(() => setLoading(false));
    }, []);

  return (
    
    // <Carousel activeIndex={index} onSelect={handleSelect}>
    loading ? (
        <div className="d-flex align-items-center">
          <strong role="status">Loading...</strong>
          <div className="spinner-border ms-auto" aria-hidden="true"></div>
        </div>
      ) : (
    <Carousel height='450' width='450' className='img-fluid bg-dark-subtle'>
        {photos.map((photo) => (
        <Carousel.Item key={photo.url}>
            <Image src={photo.url}  text="First slide" className='text-center' height='450' width='450'/>
            <Carousel.Caption key={photo.url+1000}>
            {/* <h3>First slide label</h3> */}
            <p>{photo.caption}</p>
            </Carousel.Caption>
        </Carousel.Item>
        ))}

    </Carousel>
      )
  );
}

export default ControlledCarousel;