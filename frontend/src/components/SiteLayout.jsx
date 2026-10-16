import { useEffect, useState } from 'react'
import './SiteLayout.css'

const links=[['home','Home'],['services','Services'],['about','About'],['book','Book'],['gallery','Gallery'],['contact','Contact']]

function useBranding(){
  const[branding,setBranding]=useState(null)
  useEffect(()=>{
    fetch('/config/branding.json')
      .then(response=>response.ok?response.json():Promise.reject())
      .then(config=>{
        setBranding(config)
        const favicon=document.querySelector('link[rel="icon"]')||document.createElement('link')
        favicon.rel='icon'
        favicon.href=config.favicon?.src||config.logo?.src
        if(config.favicon?.type)favicon.type=config.favicon.type
        document.head.appendChild(favicon)
        const apple=document.querySelector('link[rel="apple-touch-icon"]')||document.createElement('link')
        apple.rel='apple-touch-icon'
        apple.href=config.logo?.src
        document.head.appendChild(apple)
      })
      .catch(()=>setBranding(null))
  },[])
  return branding
}

export default function SiteLayout({page,navigate,business,children}){
  const branding=useBranding()
  const businessName=business?.project?.name||'Booksy Portal'
  const logo=branding?.logo
  return <div className="site-frame">
    <header className="site-nav">
      <button className="site-brand" onClick={()=>navigate('home')} aria-label={`${businessName} home`}>
        {logo?<img className="site-logo" src={logo.src} alt={logo.alt||`${businessName} logo`}/>:<span>BP</span>}
        {(branding?.header?.showBusinessName??true)&&<strong>{businessName}</strong>}
      </button>
      <nav aria-label="Primary navigation">{links.map(([id,label])=><button className={page===id?'active':''} onClick={()=>navigate(id)} key={id}>{label}</button>)}</nav>
      <button className="site-book" onClick={()=>navigate('book')}>Book now</button>
    </header>
    {children}
    <footer className="site-footer">
      <div className="footer-brand">
        {(branding?.footer?.showLogo??true)&&logo&&<img className="footer-logo" src={logo.src} alt=""/>}
        {(branding?.footer?.showBusinessName??true)&&<strong>{businessName}</strong>}
      </div>
      <span>{business?.business?.address||''}</span>
      <span>Powered by the configurable Booksy portal skeleton</span>
    </footer>
  </div>
}
