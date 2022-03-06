using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using l0.front.Models;
using System.Net;
using System.IO;
using System.Text.Json;

namespace l0.front.Controllers
{
    public class HomeController : Controller
    {
        private readonly ILogger<HomeController> _logger;

        public HomeController(ILogger<HomeController> logger)
        {
            _logger = logger;
        }

        public IActionResult Index()
        {
            return View();
        }

        public IActionResult GetByID()
        {
            WebRequest request = WebRequest.Create($"http://127.0.0.1/sub_cache/" + Request.Form["id_input"]);
            request.Method = "GET";

            WebResponse response = request.GetResponse();
            string respJson;
            using (Stream stream = response.GetResponseStream())
            {
                using (StreamReader reader = new StreamReader(stream))
                {
                    respJson = reader.ReadToEnd();
                }
            }

            var orderInfo = JsonSerializer.Deserialize<OrderInfo>(respJson);

            return View("Order", orderInfo);
        }

        public IActionResult Privacy()
        {
            return View();
        }

        [ResponseCache(Duration = 0, Location = ResponseCacheLocation.None, NoStore = true)]
        public IActionResult Error()
        {
            return View(new ErrorViewModel { RequestId = Activity.Current?.Id ?? HttpContext.TraceIdentifier });
        }
    }
}
