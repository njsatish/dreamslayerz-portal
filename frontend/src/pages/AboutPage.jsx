import { useEffect, useState } from 'react'

function useAbout(){
  const[data,setData]=useState(null)
  useEffect(()=>{fetch('/config/about.json').then(r=>r.ok?r.json():Promise.reject()).then(setData).catch(()=>setData(null))},[])
  return data
}

const presentationOf=provider=>provider.presentation||{}
const bookProvider=id=>{window.location.hash=`/book?provider=${encodeURIComponent(id)}`}
const openProvider=id=>{window.location.hash=`/providers/${encodeURIComponent(id)}`}

function FullProvider({provider,about,navigate}){
  const presentation=presentationOf(provider)
  const paragraphs=Array.isArray(presentation.biography)?presentation.biography:(about?.spotlight?.paragraphs||[presentation.bio||presentation.shortBio].filter(Boolean))
  const specialties=presentation.specialties||about?.spotlight?.specialties||[]
  const photo=presentation.photo||about?.spotlight?.image||'/assets/providers/provider-placeholder.jpg'
  return <section className="about-spotlight">
    <figure className="about-portrait"><img src={photo} alt={presentation.photoAlt||`${provider.name} provider portrait`}/>{(presentation.photoMode==='demo'||about?.spotlight?.imageMode==='demo')&&<figcaption>Demo provider image</figcaption>}</figure>
    <article className="about-profile-card"><p className="about-label">{about?.spotlight?.label||'Meet your provider'}</p><h2>{provider.name}</h2><h3>{provider.title||about?.spotlight?.heading}</h3>{paragraphs.map((p,i)=><p key={i}>{p}</p>)}<ul>{specialties.map(item=><li key={item}>{item}</li>)}</ul><div className="about-actions"><button onClick={()=>bookProvider(provider.id)}>Book an appointment</button><button className="outline" onClick={()=>navigate('services')}>View services</button></div></article>
    {about?.trustCard?.enabled&&<aside className="about-trust-card"><span className="trust-symbol">★</span><p>{about.trustCard.eyebrow}</p><strong>{about.trustCard.value}</strong><h3>{about.trustCard.headline}</h3><span>{about.trustCard.description}</span><button onClick={()=>bookProvider(provider.id)}>Book now</button></aside>}
  </section>
}

function TeamGrid({providers}){
  return <section className="about-team"><header><p>MEET THE TEAM</p><h2>Choose the right provider.</h2><span>Review each provider's specialties and available services, then open a profile or start booking directly.</span></header><div className="about-team-grid">{providers.map(provider=>{const p=presentationOf(provider);return <article key={provider.id}><figure><img src={p.photo||'/assets/providers/provider-placeholder.jpg'} alt={p.photoAlt||`${provider.name} provider portrait`}/>{p.photoMode==='demo'&&<figcaption>Demo image</figcaption>}</figure><div><small>{provider.title||'Service provider'}</small><h3>{provider.name}</h3><p>{p.shortBio||p.bio||'View specialties, services, and online availability.'}</p><ul>{(p.specialties||[]).slice(0,3).map(item=><li key={item}>{item}</li>)}</ul><div><button onClick={()=>openProvider(provider.id)}>View profile</button><button className="outline" onClick={()=>bookProvider(provider.id)}>Book</button></div></div></article>})}</div></section>
}

export default function AboutPage({config,navigate}){
  const about=useAbout()
  const providers=(config?.providers?.providers||[]).filter(provider=>provider.enabled).sort((a,b)=>(a.presentation?.displayOrder||0)-(b.presentation?.displayOrder||0))
  if(!about)return <main className="about-loading"><p>Loading About page…</p></main>
  return <><section className="about-editorial-hero"><div><p>{providers.length>1?'About the team':about.eyebrow}</p><h1>{providers.length>1?'Experienced providers. Clear booking.':about.headline}</h1><span>{providers.length>1?'Meet each professional, review specialties, and choose the right provider for your appointment.':about.introduction}</span></div></section>{providers.length===1?<FullProvider provider={providers[0]} about={about} navigate={navigate}/>:providers.length>1?<TeamGrid providers={providers}/>:<section className="about-spotlight"><article className="about-profile-card"><p className="about-label">ABOUT THE BUSINESS</p><h2>Provider profiles are coming soon.</h2><p>Business information and online booking remain available.</p><div className="about-actions"><button onClick={()=>navigate('book')}>Book an appointment</button></div></article></section>}<section className="about-values"><header><p>THE EXPERIENCE</p><h2>Built around clarity and convenience.</h2></header><div className="card-grid"><article className="content-card"><span>01</span><h3>Choose</h3><p>Review services, durations, prices, and eligible providers.</p></article><article className="content-card"><span>02</span><h3>Schedule</h3><p>Select current dates and times from live availability.</p></article><article className="content-card"><span>03</span><h3>Confirm</h3><p>Complete the appointment securely through Booksy.</p></article></div></section></>
}
