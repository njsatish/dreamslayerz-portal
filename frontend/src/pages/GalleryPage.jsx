import { useEffect, useMemo, useState } from 'react'
import './GalleryPage.css'

function selectedImageFromHash(){
  const query=window.location.hash.split('?')[1]||''
  return new URLSearchParams(query).get('image')||''
}

function useGallery(){
  const [data,setData]=useState(null)
  useEffect(()=>{fetch('/config/gallery.json').then(r=>r.ok?r.json():Promise.reject()).then(setData).catch(()=>setData({enabled:false,items:[]}))},[])
  return data
}

export default function GalleryPage({onBack}){
  const data=useGallery()
  const [category,setCategory]=useState('All')
  const [activeIndex,setActiveIndex]=useState(null)
  const [requestedImage,setRequestedImage]=useState(selectedImageFromHash())
  const categories=useMemo(()=>['All',...new Set((data?.items||[]).map(item=>item.category))],[data])
  const items=useMemo(()=>category==='All'?(data?.items||[]):(data?.items||[]).filter(item=>item.category===category),[data,category])

  useEffect(()=>{setActiveIndex(null)},[category])
  useEffect(()=>{
    if(!data||!requestedImage)return
    const index=(data.items||[]).findIndex(item=>item.id===requestedImage)
    if(index>=0){
      setCategory('All')
      setActiveIndex(index)
    }
    setRequestedImage('')
  },[data,requestedImage])
  useEffect(()=>{
    if(activeIndex===null)return
    const previous=document.body.style.overflow
    document.body.style.overflow='hidden'
    const keydown=event=>{
      if(event.key==='Escape')setActiveIndex(null)
      if(event.key==='ArrowRight')setActiveIndex(index=>(index+1)%items.length)
      if(event.key==='ArrowLeft')setActiveIndex(index=>(index-1+items.length)%items.length)
    }
    window.addEventListener('keydown',keydown)
    return()=>{document.body.style.overflow=previous;window.removeEventListener('keydown',keydown)}
  },[activeIndex,items.length])

  if(!data)return <main className="gallery-page"><p>Loading gallery…</p></main>
  if(!data.enabled)return <main className="gallery-page"><button className="gallery-back" onClick={onBack}>← Back to home</button><h1>Gallery unavailable</h1></main>
  const active=activeIndex===null?null:items[activeIndex]
  return <main className="gallery-page">
    <header className="gallery-page-header">
      <button className="gallery-back" onClick={onBack}>← Back to home</button>
      <p>PORTFOLIO</p><h1>{data.heading}</h1><span>{data.intro}</span>
      {data.mode==='demo'&&<div className="gallery-disclosure"><strong>Demo content:</strong> {data.disclosure}</div>}
    </header>
    <nav className="gallery-filters" aria-label="Gallery categories">
      {categories.map(name=><button className={category===name?'active':''} onClick={()=>setCategory(name)} key={name}>{name}</button>)}
    </nav>
    <div className="gallery-page-grid">{items.map((item,index)=><button className="gallery-tile" onClick={()=>setActiveIndex(index)} key={item.id}>
      <img src={item.src} alt={item.alt} loading="lazy" width="700" height="700"/>
      <span className="gallery-tile-overlay"><strong>{item.title||item.category}</strong><small>{item.generated?'AI placeholder':'Business portfolio'}</small></span>
    </button>)}</div>
    {!items.length&&<p className="gallery-empty">No images are configured for this category.</p>}
    {active&&<div className="gallery-viewer" role="dialog" aria-modal="true" aria-label={active.title||active.category}>
      <button className="gallery-viewer-backdrop" onClick={()=>{setActiveIndex(null);window.location.hash='/gallery'}} aria-label="Close image viewer"/>
      <div className="gallery-viewer-dialog">
        <button className="gallery-viewer-close" onClick={()=>{setActiveIndex(null);window.location.hash='/gallery'}} aria-label="Close image viewer">×</button>
        <button className="gallery-viewer-prev" onClick={()=>setActiveIndex((activeIndex-1+items.length)%items.length)} aria-label="Previous image">‹</button>
        <img src={active.src} alt={active.alt}/>
        <button className="gallery-viewer-next" onClick={()=>setActiveIndex((activeIndex+1)%items.length)} aria-label="Next image">›</button>
        <footer><strong>{active.title||active.category}</strong><span>{active.generated?'AI-generated layout placeholder':'Business portfolio image'}</span><small>{activeIndex+1} of {items.length}</small></footer>
      </div>
    </div>}
  </main>
}
