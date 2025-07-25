import { useState, useEffect } from 'react';
import Carousel from 'react-bootstrap/Carousel';
import Image from 'react-bootstrap/Image';

function ControlledCarousel({photos, id}) {
//   const [index, setIndex] = useState(0);

//   const handleSelect = (selectedIndex) => {
//     setIndex(selectedIndex);
//   };

    // const [photos, setPhotos] = useState([]);
    // const [loading, setLoading] = useState(true);
  
    // useEffect(() => {
    //   fetch(`http://localhost:8080/photos/${id}`)
    //     .then((res) => res.json())
    //     .then((data) => {
    //       setPhotos(data);
    //       setLoading(false);
    //     })
    //     .catch(() => setLoading(false));
    // }, []);

  return (
    
    // <Carousel activeIndex={index} onSelect={handleSelect}>
    photos==undefined || photos.length==0 ? (
        <div className="d-flex align-items-center">
          <strong role="status">Loading...</strong>
          <div className="spinner-border ms-auto text-center" aria-hidden="true"></div>
          <p>This user has no images to display.</p>
        </div>
      ) : (
    <Carousel height='450' width='450' className='img-fluid bg-dark-subtle'>
        {photos.map((photo) => (
        <Carousel.Item key={photo.url+id}>
            <Image src={photo.url}  text="First slide" className='text-center' style={{ maxWidth: 'auto', maxHeight: '500px'}}/>
            <Carousel.Caption key={photo.url+1000}>
            {/* <h3>First slide label</h3> */}
            <span className='bg-secondary m-2 p-2 rounded' style={{ '--bs-bg-opacity': '.7' }}>{photo.caption}</span>
            </Carousel.Caption>
        </Carousel.Item>
        ))}

    </Carousel>
      )
  );
}

export default ControlledCarousel;