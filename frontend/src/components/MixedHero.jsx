import { useEffect, useMemo, useRef, useState } from 'react'
import './MixedHero.css'

function useJSON(path){
  const[data,setData]=useState(null)
  useEffect(()=>{fetch(path).then(r=>r.ok?r.json():Promise.reject()).then(setData).catch(()=>setData(null))},[path])
  return data
}

const duration=minutes=>minutes<60?`${minutes} min`:`${Math.floor(minutes/60)} hr${minutes%60?` ${minutes%60} min`:''}`

export default function MixedHero({config,navigate,openGallery}){
  const hero=useJSON('/config/hero.json')
  const[index,setIndex]=useState(0)
  const[paused,setPaused]=useState(false)
  const rootRef=useRef(null)
  const slides=hero?.slides||[]
  const active=slides[index]||null

  const details=useMemo(()=>{
    if(!active||!config)return{}
    if(active.type==='service'){
      const service=config.services.services.find(item=>item.slug===active.serviceSlug)
      const variant=service?.variants.find(item=>item.id===active.variantId)||service?.variants[0]
      return {price:variant?.priceLabel,duration:variant?.durationMinutes?duration(variant.durationMinutes):'',service,variant}
    }
    if(active.type==='provider'){
      const provider=config.providers.providers.find(item=>item.id===active.providerId)
      return {provider}
    }
    return{}
  },[active,config])

  const move=direction=>setIndex(current=>(current+direction+slides.length)%slides.length)
  useEffect(()=>{
    if(!hero?.carousel?.autoplay||paused||slides.length<2)return
    const timer=window.setInterval(()=>move(1),Math.max(3,hero.carousel.intervalSeconds||7)*1000)
    return()=>window.clearInterval(timer)
  },[hero,paused,slides.length])
  useEffect(()=>{
    const element=rootRef.current
    if(!element)return
    const key=event=>{if(event.key==='ArrowRight')move(1);if(event.key==='ArrowLeft')move(-1)}
    element.addEventListener('keydown',key)
    return()=>element.removeEventListener('keydown',key)
  },[slides.length])
  if(!hero?.enabled||!active)return null

  const act=slide=>{
    if(slide.action?.target==='gallery')return openGallery(slide.galleryItemId)
    if(slide.action?.target==='book'){
      const query=new URLSearchParams()
      if(slide.serviceSlug)query.set('service',slide.serviceSlug)
      if(slide.providerId)query.set('provider',slide.providerId)
      window.location.hash=`/book${query.size?`?${query}`:''}`
      return
    }
    navigate(slide.action?.target||'home')
  }

  return <section className="mixed-hero" ref={rootRef} tabIndex="-1" onMouseEnter={()=>hero.carousel.pauseOnHover&&setPaused(true)} onMouseLeave={()=>setPaused(false)}>
    <div className="mixed-hero-copy">
      <p className="mixed-eyebrow">{hero.eyebrow}</p>
      <h1><span>{hero.headline.line1}</span><strong>{hero.headline.line2}</strong></h1>
      <p className="mixed-description">{hero.description}</p>
      <div className="mixed-actions">
        <button className="mixed-primary" onClick={()=>navigate(hero.primaryAction.target)}>{hero.primaryAction.label}</button>
        <button className="mixed-secondary" onClick={()=>navigate(hero.secondaryAction.target)}>{hero.secondaryAction.label}</button>
      </div>
      <div className="mixed-proof"><span>LIVE</span><p>Service-led booking for one provider or many.</p></div>
    </div>

    <div className="mixed-spotlight" aria-live="polite">
      <img src={active.image} alt=""/>
      <div className="mixed-shade"/>
      <div className="mixed-badge"><small>{active.badge}</small><strong>{active.title}</strong></div>
      {active.type==='service'&&<div className="mixed-stat"><small>Starting at</small><strong>{details.price||'See service'}</strong><span>{details.duration||'Duration varies'}</span></div>}
      {active.type==='provider'&&<div className="mixed-stat"><small>Provider</small><strong>{details.provider?.name||active.title}</strong><span>{details.provider?.title||'Service provider'}</span></div>}
      {(active.type==='gallery'||active.type==='booking')&&<div className="mixed-stat"><small>{active.type==='gallery'?'Portfolio':'Booking'}</small><strong>{active.type==='gallery'?'View work':'Live times'}</strong><span>{active.type==='gallery'?'Open exact image':'Powered by Booksy'}</span></div>}
      <div className="mixed-caption"><div><small>{active.type}</small><strong>{active.title}</strong><span>{active.subtitle}</span></div><button onClick={()=>act(active)}>{active.action.label} →</button></div>
      {hero.carousel.showArrows&&slides.length>1&&<><button className="mixed-arrow previous" onClick={()=>move(-1)} aria-label="Previous hero slide">‹</button><button className="mixed-arrow next" onClick={()=>move(1)} aria-label="Next hero slide">›</button></>}
      {hero.carousel.showIndicators&&<div className="mixed-indicators" aria-label="Hero slides">{slides.map((slide,slideIndex)=><button className={slideIndex===index?'active':''} onClick={()=>setIndex(slideIndex)} aria-label={`Show slide ${slideIndex+1}: ${slide.title}`} key={slide.id}/>)}</div>}
    </div>
  </section>
}
