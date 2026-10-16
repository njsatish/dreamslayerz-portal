import { GallerySection, ReviewsSection } from '../content/ContentSections'
import MixedHero from '../components/MixedHero'

export default function HomePage({config,navigate,openGallery}){
  const hero=config?.theme?.hero
  const address=config?.business?.business?.address
  return <>
    <MixedHero config={config} navigate={navigate} openGallery={openGallery}/>

    <section className="home-booking-cta">
      <div><p>LIVE AVAILABILITY</p><h2>Ready when you are.</h2><span>Choose a service, select an available date and time, then complete the appointment securely through Booksy.</span></div>
      <button onClick={()=>navigate('book')}>Find an appointment</button>
    </section>

    <GallerySection onOpenGallery={openGallery}/>
    <ReviewsSection/>
  </>
}
