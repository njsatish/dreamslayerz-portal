import './ProviderPage.css'

function providerServices(provider,services){
  return (services||[]).map(service=>{
    const variants=(service.variants||[]).filter(variant=>(variant.providerIds||[]).includes(provider.id))
    return variants.length?{...service,variants}:null
  }).filter(Boolean)
}

function bookProvider(providerId,serviceSlug){
  const query=new URLSearchParams({provider:providerId})
  if(serviceSlug)query.set('service',serviceSlug)
  window.location.hash=`/book?${query}`
}

export default function ProviderPage({config,providerId,navigate}){
  const routeProviderId=providerId||decodeURIComponent((window.location.hash.match(/^#?\/providers\/([^?]+)/)||[])[1]||'')
  const providers=(config?.providers?.providers||[]).filter(provider=>provider.enabled)
  const provider=providers.find(item=>item.id===routeProviderId)
  if(!provider)return <main className="provider-not-found">
    <p>PROVIDER DIRECTORY</p>
    <h1>Provider not found.</h1>
    <span>The requested provider is unavailable or is no longer publicly listed.</span>
    <div><button onClick={()=>navigate('about')}>Meet the team</button><button className="outline" onClick={()=>navigate('book')}>Book an appointment</button></div>
  </main>

  const presentation=provider.presentation||{}
  const photo=presentation.photo||'/assets/providers/provider-placeholder.jpg'
  const specialties=presentation.specialties||[]
  const biography=Array.isArray(presentation.biography)?presentation.biography:[presentation.bio||presentation.shortBio].filter(Boolean)
  const offerings=providerServices(provider,config?.services?.services)

  return <>
    <section className="provider-profile-hero">
      <figure><img src={photo} alt={presentation.photoAlt||`${provider.name} provider portrait`}/>{presentation.photoMode==='demo'&&<figcaption>Demo provider image</figcaption>}</figure>
      <div className="provider-profile-copy">
        <p>{provider.title||'Service provider'}</p>
        <h1>{provider.name}</h1>
        <span>{presentation.shortBio||'Review specialties, services, and current booking options.'}</span>
        <div className="provider-profile-actions"><button onClick={()=>bookProvider(provider.id)}>Book with {provider.name}</button><button className="outline" onClick={()=>navigate('about')}>Back to team</button></div>
      </div>
    </section>

    <section className="provider-profile-body">
      <article className="provider-biography">
        <p className="provider-label">ABOUT THE PROVIDER</p>
        <h2>Experience and approach</h2>
        {biography.length?biography.map((paragraph,index)=><p key={index}>{paragraph}</p>):<p>Add the approved provider biography in <code>config/providers.json</code>.</p>}
        {!!specialties.length&&<><h3>Specialties</h3><ul>{specialties.map(item=><li key={item}>{item}</li>)}</ul></>}
      </article>
      <aside className="provider-quick-book"><p>READY TO BOOK?</p><h2>Choose this provider.</h2><span>Select an eligible service, live date, and available time.</span><button onClick={()=>bookProvider(provider.id)}>View availability</button></aside>
    </section>

    <section className="provider-services-section">
      <header><p>AVAILABLE SERVICES</p><h2>Book with {provider.name}</h2></header>
      {offerings.length?<div className="provider-services-grid">{offerings.map(service=><article key={service.slug}>
        <small>{service.category||'Service'}</small><h3>{service.name}</h3>
        <div>{service.variants.map(variant=><p key={variant.id}><span>{variant.label||service.name}</span><strong>{variant.priceLabel||''}</strong><small>{variant.durationMinutes?`${variant.durationMinutes} min`:''}</small></p>)}</div>
        <button onClick={()=>bookProvider(provider.id,service.slug)}>Book this service →</button>
      </article>)}</div>:<div className="provider-no-services">No online services are currently assigned to this provider.</div>}
    </section>
  </>
}
