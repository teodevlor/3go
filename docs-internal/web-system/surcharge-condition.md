// 1. template time_window

{
  "code": "RUSH_HOUR_17_19_WEEKDAY",
  "condition_type": "time_window",
  "config": {
    "from": "17:00",
    "to": "19:00",
    "days": ["mon", "tue", "wed", "thu", "fri"]
  },
  "is_active": true
}

// 2.1 template weather begin
{
  "code": "WEATHER_RAIN",
  "condition_type": "weather",
  "config": {
    "rain": true
  },
  "is_active": true
}

// 2.2 template weather >= 10mm
{
  "code": "HEAVY_RAIN_OVER_10MM",
  "condition_type": "weather",
  "config": {
    "rain_mm": {
      "operator": ">=",
      "value": 10
    }
  },
  "is_active": true
}

// 3. template traffic

{
  "code": "TRAFFIC_HIGH",
  "condition_type": "traffic",
  "config": {
    "level": "high"
  },
  "is_active": true
}

// 4. template holiday

{
  "code": "HOLIDAY_TET",
  "condition_type": "holiday",
  "config": {
    "holiday_code": "TET"
  },
  "is_active": true
}

