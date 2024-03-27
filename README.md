How would you scale this script and run it with resiliency to e.g. handle 1000s of domains?
- I would not use the public API, but rather stand up an instance myself or compile the CLI locally,
- Use goroutines and channels to handle parallelism

How would you monitor/alert on this service?
- A simple solution would be to export Prometheus metrics to be scraped by something like node exporter if this were running in Kubernetes
- In the code itself if an error or failed scan reach out to some external alert service. 

What would you do to handle adding new domains to scan or certificate expiry events from your service?
- Scrape the list of domains to scan from some external endpoint. 
- When a cert expire event is near or occurs, send out an email alert or hook into some other alert service like PagerDuty 

After some time, your report requires more enhancements requested by the Tech team of the company. How would you handle these "continuous" requirement changes in a sustainable manner?
- Short answer: gitops
- establish a ci/cd pipeline for the application stack