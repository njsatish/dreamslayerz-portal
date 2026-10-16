import { useCallback, useEffect, useRef, useState } from 'react'
import './ContentSections.css'
import HomeGalleryLightbox from './HomeGalleryLightbox'

function useContent(path){
  const [data,setData]=useState(null)
  useEffect(()=>{fetch(path).then(r=>r.ok?r.json():Promise.reject()).then(setData).catch(()=>setData({enabled:false}))},[path])
  return data
}

function Disclosure({children}){return <p className="demo-disclosure"><strong>Demo content:</strong> {children}</p>}

export function GallerySection({onOpenGallery}){
  const [lightboxIndex,setLightboxIndex]=useState(null)
  const triggerRef=useRef(null)
  const closePreview=useCallback(()=>{
    setLightboxIndex(null)
    requestAnimationFrame(()=>triggerRef.current?.focus())
  },[])

  const data=useContent('/config/gallery.json')
  if(!data?.enabled)return null
  return <section className="content-section gallery-section" id="gallery">{lightboxIndex!==null&&<HomeGalleryLightbox items={data.items.slice(0,4)} index={lightboxIndex} onIndexChange={setLightboxIndex} onClose={closePreview} onViewGallery={selected=>{closePreview();onOpenGallery(selected.id||selected)}}/>}
    <header><p>PORTFOLIO</p><h2>{data.heading}</h2><span>{data.intro}</span></header>
    {data.mode==='demo'&&<Disclosure>{data.disclosure}</Disclosure>}
    <div className="gallery-grid">{data.items.slice(0,4).map((item,index)=><button className="gallery-preview-card" onClick={event=>{triggerRef.current=event.currentTarget;setLightboxIndex(index)}} role="button" tabIndex={0} onKeyDown={event=>{if(event.key==='Enter'||event.key===' '){event.preventDefault();triggerRef.current=event.currentTarget;setLightboxIndex(index)}}} key={item.id}>
      <img src={item.src} alt={item.alt} loading="lazy" width="600" height="600"/>
      <span className="gallery-preview-caption"><strong>{item.title||item.category}</strong>{item.generated&&<small>AI placeholder</small>}</span>
    </button>)}</div><button className="view-gallery-button" onClick={onOpenGallery}>View full gallery →</button>
  </section>
}

export function ReviewsSection(){
  const data=useContent('/config/reviews.json')
  if(!data?.enabled)return null
  return <section className="content-section reviews-section" id="reviews">
    <header><p>REVIEWS</p><h2>{data.heading}</h2><span>{data.intro}</span></header>
    {data.mode==='demo'&&<Disclosure>{data.disclosure}</Disclosure>}
    <div className="reviews-grid">{data.reviews.map(review=><article key={review.id}>
      <div className="stars" aria-label={`${review.rating} out of 5 stars`}>{'★'.repeat(review.rating)}</div>
      <blockquote>“{review.text}”</blockquote>
      <footer><strong>{review.displayName}</strong><span>{review.verified?'Verified review':'Sample review'}</span></footer>
    </article>)}</div>
  </section>
}
