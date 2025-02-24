const puppeteer = require('puppeteer');
const fs = require('fs');
const path = require('path');

(async () => {
  try {
    console.log("🚀 Iniciando Puppeteer para obtener la cookie de Vinted...");

    
    const browser = await puppeteer.launch({ headless: "new" });
    const page = await browser.newPage();

    
    await page.goto('https://www.vinted.es', { waitUntil: 'networkidle2' });

    
    await new Promise(resolve => setTimeout(resolve, 3000));

    
    const cookies = await page.cookies();

   
    const sessionCookie = cookies.find(c => c.name === '_vinted_fr_session');

    if (!sessionCookie) {
      console.log('❌ No encontré la cookie _vinted_fr_session.');
      process.exit(1);
    }

    console.log('✅ Cookie encontrada:');
    console.log(`Name: ${sessionCookie.name}`);
    console.log(`Value: ${sessionCookie.value}`);

    
    const cookiePath = path.join(__dirname, '../snoopycoon_back/vinted_cookie.txt');
    fs.writeFileSync(cookiePath, sessionCookie.value);
    console.log(`✅ La cookie se guardó en: ${cookiePath}`);

    
    await browser.close();
    process.exit(0);
  } catch (err) {
    console.error('❌ Error en getVintedCookie:', err);
    process.exit(1);
  }
})();




