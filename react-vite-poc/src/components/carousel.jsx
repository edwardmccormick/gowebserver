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
    <Carousel height='450' width='450' className='ratio ratio-1x1'>
        {photos.map((photo) => (
        <Carousel.Item >
            <Image src={photo.url}  text="First slide" className='text-center' height='450' width='450'/>
            <Carousel.Caption>
            {/* <h3>First slide label</h3> */}
            <p>{photo.caption}</p>
            </Carousel.Caption>
        </Carousel.Item>
        ))}

    </Carousel>
  );
}

export default ControlledCarousel;