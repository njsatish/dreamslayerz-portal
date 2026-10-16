import { useEffect, useRef } from 'react'
import './HomeGalleryLightbox.css'

export default function HomeGalleryLightbox({items,index,onIndexChange,onClose,onViewGallery}){
  const closeRef=useRef(null)
  const active=items[index]

  useEffect(()=>{
    const previousOverflow=document.body.style.overflow
    document.body.style.overflow='hidden'
    closeRef.current?.focus()
    const onKeyDown=event=>{
      if(event.key==='Escape')onClose()
      if(event.key==='ArrowLeft')onIndexChange((index-1+items.length)%items.length)
      if(event.key==='ArrowRight')onIndexChange((index+1)%items.length)
      if(event.key==='Tab'){
        const controls=[...document.querySelectorAll('.home-lightbox button')]
        if(!controls.length)return
        const first=controls[0],last=controls[controls.length-1]
        if(event.shiftKey&&document.activeElement===first){event.preventDefault();last.focus()}
        else if(!event.shiftKey&&document.activeElement===last){event.preventDefault();first.focus()}
      }
    }
    window.addEventListener('keydown',onKeyDown)
    return()=>{document.body.style.overflow=previousOverflow;window.removeEventListener('keydown',onKeyDown)}
  },[index,items.length,onClose,onIndexChange])

  if(!active)return null
  const src=active.src||active.image||active.url
  const title=active.title||active.name||`Work photo ${index+1}`
  return <div className="home-lightbox" role="dialog" aria-modal="true" aria-label={`Work photo preview: ${title}`} onMouseDown={event=>event.target===event.currentTarget&&onClose()}>
    <div className="home-lightbox-dialog">
      <button ref={closeRef} className="home-lightbox-close" onClick={onClose} aria-label="Close photo preview">×</button>
      <button className="home-lightbox-arrow previous" onClick={()=>onIndexChange((index-1+items.length)%items.length)} aria-label="Previous work photo">‹</button>
      <figure>
        <img src={src} alt={active.alt||title}/>
        <figcaption>
          <div><small>{index+1} of {items.length}</small><strong>{title}</strong>{active.description&&<span>{active.description}</span>}</div>
          <button onClick={()=>onViewGallery(active)}>View full gallery →</button>
        </figcaption>
      </figure>
      <button className="home-lightbox-arrow next" onClick={()=>onIndexChange((index+1)%items.length)} aria-label="Next work photo">›</button>
    </div>
  </div>
}
